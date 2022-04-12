package stonks

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"text/template"
	"time"

	cg "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
	"go.uber.org/ratelimit"
)

type Stonk interface {
	SetCurrency(c string)

	GetGems() ([]Coin, error)

	GetCoinList() ([]Coin, error)
	GetPrice(id string) float32
	GetCoinDataFromID(id string) (CoinData, error)

	ToHTML(w io.Writer, data interface{}) error
}

type Coin struct {
	ID     string
	Symbol string
	Name   string
}

type CoinData struct {
	Coin
	MarketCapRank       uint16
	MarketCap           float64
	DeveloperScore      float32
	PublicInterestScore float32
	CommunityScore      float32
	ImageURL            string
	Description         string
	CurrentPrice        float64

	MoonshotScore float32
}

func (cd *CoinData) score() {
	cd.MoonshotScore = 2*cd.DeveloperScore + cd.PublicInterestScore + 100*float32(cd.MarketCapRank)
}

type Stonker struct {
	client          *cg.Client
	ratelimiter     ratelimit.Limiter
	defaultCurrency string
	web             WebData
}

type WebData struct {
	Template  string
	Path      string
	Gems      []CoinData
	Timestamp string
}

func NewStonker() *Stonker {
	httpClient := &http.Client{
		Timeout: time.Second * 90,
	}
	return &Stonker{
		client:          cg.NewClient(httpClient),
		ratelimiter:     ratelimit.New(1, ratelimit.Per(2*time.Second)), // 1 request per 2 seconds (0.5 rps)
		defaultCurrency: "aud",
		web: WebData{
			Template:  indexTemplate,
			Path:      "./web/index.html",
			Gems:      []CoinData{},
			Timestamp: time.Now().Format(time.RFC850),
		},
	}
}

func getData(s Stonker, coins []CoinData, coin Coin, idx int, wg *sync.WaitGroup) CoinData {
	var cd CoinData = CoinData{}
	var err error = nil
	if coin.ID != "" {
		cd, err = s.GetCoinDataFromID(coin.ID)
		if err != nil {
			log.Printf("  - ERROR fetching '%+v': %s\n", coin, err)
		}
	}
	wg.Done()
	coins[idx] = cd
	return cd
}

func (s Stonker) GetGems(top int) ([]CoinData, error) {
	list, err := s.GetCoinList()
	if err != nil {
		return nil, err
	}

	// shuffle for fairness in case of failures
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	log.Printf("found %d potential coins.\n", len(list))

	// TODO(scharissis): remove :)
	// temporary hack to allow this to complete within 6hrs on Github Actions
	list = list[:10000]

	var wg sync.WaitGroup
	coins := make([]CoinData, len(list))
	log.Println(" - requesting coin data...")
	for i, c := range list {
		wg.Add(1)
		go getData(s, coins, c, i, &wg)
	}
	log.Println(" - waiting for coin data...")
	wg.Wait()
	logSuccessRate(list, coins)

	coins = rankAndFilter(coins)
	log.Printf("found %d potential gems.\n", len(coins))
	coins = coins[:top]

	log.Printf("Top %d gems:\n", top)
	for i, c := range coins {
		log.Printf("#%d: %+v\n", i+1, c)
	}
	return coins, err
}

func rankAndFilter(coins []CoinData) []CoinData {
	coins = filter(coins,
		badDeveloperScore,
		topTenMarketCap,
	)
	coins = rank(coins)
	return coins
}

func (s Stonker) GetCoinList() ([]Coin, error) {
	s.ratelimiter.Take() // rate limit
	coinList, err := s.client.CoinsList()
	if err != nil {
		return nil, err
	}
	coins := []Coin{}
	for _, c := range *coinList {
		coins = append(coins, Coin{ID: c.ID, Name: c.Name, Symbol: c.Symbol})
	}
	return coins, err
}

func (s Stonker) GetPrice(id string) float32 {
	s.ratelimiter.Take() // rate limit
	price, err := s.client.SimpleSinglePrice(id, s.defaultCurrency)
	if err != nil {
		log.Fatal(err)
	}
	return price.MarketPrice
}

func (s Stonker) GetCoinDataFromID(id string) (CoinData, error) {
	s.ratelimiter.Take() // rate limit
	coin, err := s.client.CoinsID(id, false, true, true, true, true, true)
	if err != nil {
		return CoinData{}, err
	}
	return CoinDataFromCoinsID(coin, s.defaultCurrency), err
}

func CoinDataFromCoinsID(cid *types.CoinsID, currency string) CoinData {
	description := "no description found"
	if d, found := cid.Description["en"]; found {
		d = strings.Split(d, `.`)[0]
		description = d
	}
	currPrice := 0.0
	if p, found := cid.MarketData.CurrentPrice[currency]; found {
		currPrice = p
	}
	cd := CoinData{
		Coin: Coin{
			ID:     cid.ID,
			Name:   cid.Name,
			Symbol: cid.Symbol,
		},
		MarketCapRank:       cid.MarketCapRank,
		MarketCap:           cid.MarketData.MarketCap[currency],
		DeveloperScore:      cid.DeveloperScore,
		PublicInterestScore: cid.PublicInterestScore,
		CommunityScore:      cid.CommunityScore,
		ImageURL:            cid.Image.Large,
		Description:         description,
		CurrentPrice:        currPrice,
	}
	cd.score()
	return cd
}

func (s Stonker) ToHTML(w io.Writer) error {
	log.Printf("generating html...\n")
	var err error = nil
	s.web.Gems, err = s.GetGems(3)
	if err != nil {
		return err
	}

	log.Printf(" - with %d gems\n", len(s.web.Gems))
	funcMap := template.FuncMap{
		"GetCurrency": func() string { return s.defaultCurrency },
		"ToUpper":     strings.ToUpper,
	}
	t := template.Must(template.New("html").Funcs(funcMap).Parse(s.web.Template))
	return t.Execute(w, s.web)
}

func logSuccessRate(list []Coin, coins []CoinData) {
	all := len(list)
	success := 0
	emptyCoinData := CoinData{}
	for _, c := range coins {
		if !reflect.DeepEqual(c, emptyCoinData) {
			success++
		}
	}
	log.Printf("%d/%d (%.1f%%) coins successfully queried.\n", success, all, float32(success)/float32(all)*100.0)
}

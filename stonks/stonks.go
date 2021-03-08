package stonks

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cg "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type Stonk interface {
	SetCurrency(c string)

	GetGems() ([]Coin, error)

	GetCoinList() ([]Coin, error)
	GetPrice(id string) float32
	GetCoinDataFromID(id string) (CoinData, error)
}

type Coin struct {
	ID     string
	Symbol string
	Name   string
}

type CoinData struct {
	Coin
	MarketCapRank  uint16
	DeveloperScore float32
}

type Stonker struct {
	client          *cg.Client
	defaultCurrency string
}

func NewStonker() *Stonker {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Stonker{
		client:          cg.NewClient(httpClient),
		defaultCurrency: "aud",
	}
}

func getData(s Stonker, coin Coin) CoinData {
	cd, err := s.GetCoinDataFromID(coin.ID)
	if err != nil {
		fmt.Printf("  - ERROR fetching %+s: %s\n", coin.ID, err)
		return CoinData{}
	}
	return cd
}

func (s Stonker) GetGems() ([]CoinData, error) {
	list, err := s.GetCoinList()
	if err != nil {
		return nil, err
	}
	coins := []CoinData{}

	for _, c := range list {
		cd := getData(s, c)
		coins = append(coins, cd)
		time.Sleep(300 * time.Millisecond)
	}

	coins = filter(coins,
		func(cd CoinData) bool { return cd.DeveloperScore == 0 },
	)
	fmt.Printf("found %d gems out of %d coins.\n", len(coins), len(list))

	rank(coins)
	for i, c := range coins[:10] {
		fmt.Printf("#%d: %+v\n", i+1, c)
	}
	return coins, err
}

func (s Stonker) GetCoinList() ([]Coin, error) {
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
	price, err := s.client.SimpleSinglePrice(id, s.defaultCurrency)
	if err != nil {
		log.Fatal(err)
	}
	return price.MarketPrice
}

func (s Stonker) GetCoinDataFromID(id string) (CoinData, error) {
	coin, err := s.client.CoinsID(id, false, true, true, true, true, true)
	if err != nil {
		return CoinData{}, err
	}
	return CoinDataFromCoinsID(coin), err
}

func CoinDataFromCoinsID(cid *types.CoinsID) CoinData {
	return CoinData{
		Coin: Coin{
			ID:     cid.ID,
			Name:   cid.Name,
			Symbol: cid.Symbol,
		},
		MarketCapRank:  cid.MarketCapRank,
		DeveloperScore: cid.DeveloperScore,
	}
}

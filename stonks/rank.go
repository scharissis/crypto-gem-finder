package stonks

import (
	"sort"
)

func rank(cds []CoinData) []CoinData {
	ranked := make([]CoinData, len(cds))
	for i, cd := range cds {
		cd.score()
		ranked[i] = cd
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].MoonshotScore > ranked[j].MoonshotScore
	})
	return ranked
}

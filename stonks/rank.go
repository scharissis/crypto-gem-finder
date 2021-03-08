package stonks

import "sort"

func rank(cds []CoinData) {
	sort.Slice(cds, func(i, j int) bool {
		return cds[i].DeveloperScore > cds[j].DeveloperScore
	})
}

package stonks

// FilterFn removes a CoinData if it returns true.
type FilterFn func(cd CoinData) bool

var badDeveloperScore FilterFn = func(cd CoinData) bool {
	return cd.DeveloperScore == 0
}

var topTenMarketCap FilterFn = func(cd CoinData) bool {
	return cd.MarketCapRank <= 10
}

// Remove element for FilterFn which are satisfied.
func filter(cds []CoinData, fns ...FilterFn) []CoinData {
	fcds := []CoinData{}
	for _, cd := range cds {
		keep := true
		for _, fn := range fns {
			if fn(cd) {
				keep = false
				break
			}
		}
		if keep {
			fcds = append(fcds, cd)
		}
	}
	return fcds
}

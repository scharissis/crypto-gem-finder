package stonks

type FilterFn func(cd CoinData) bool

var badDeveloperScore FilterFn = func(cd CoinData) bool {
	return cd.DeveloperScore == 0
}

// Remove element if any FilterFn returns true.
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

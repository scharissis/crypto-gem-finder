package stonks

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_filter(t *testing.T) {
	tests := []struct {
		cds      []CoinData
		fns      []FilterFn
		expected []CoinData
	}{
		{
			cds: []CoinData{
				{Coin: Coin{ID: "a"}, DeveloperScore: 13.37},
				{Coin: Coin{ID: "b"}, DeveloperScore: 0.0},
			},
			fns: []FilterFn{
				badDeveloperScore,
			},
			expected: []CoinData{
				{Coin: Coin{ID: "a"}, DeveloperScore: 13.37},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := filter(test.cds, test.fns...)
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("got %v, want %v\n", got, test.expected)
			}
		})
	}
}

func Test_score(t *testing.T) {
	// cd.MoonshotScore = 2*cd.DeveloperScore + cd.PublicInterestScore + 100*float32(cd.MarketCapRank)
	tests := []struct {
		cd       CoinData
		expected float32
	}{
		{cd: CoinData{Coin: Coin{ID: "b"}, DeveloperScore: 10, PublicInterestScore: 5, MarketCapRank: 3}, expected: 325},
		{cd: CoinData{Coin: Coin{ID: "c"}, DeveloperScore: 50, PublicInterestScore: 15, MarketCapRank: 0}, expected: 115},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			test.cd.score()
			if test.cd.MoonshotScore != test.expected {
				t.Fatalf("got %v, want %v\n", test.cd.MoonshotScore, test.expected)
			}
		})
	}
}

func Test_rank(t *testing.T) {
	tests := []struct {
		cds      []CoinData
		expected []CoinData
	}{
		{
			cds: []CoinData{
				{Coin: Coin{ID: "a"}, DeveloperScore: 10, PublicInterestScore: 5, MarketCapRank: 3},
				{Coin: Coin{ID: "b"}, DeveloperScore: 50, PublicInterestScore: 15, MarketCapRank: 5},
				{Coin: Coin{ID: "c"}, DeveloperScore: 40, PublicInterestScore: 10, MarketCapRank: 4},
			},
			expected: []CoinData{
				{Coin: Coin{ID: "b"}, DeveloperScore: 50, PublicInterestScore: 15, MarketCapRank: 5, MoonshotScore: 615},
				{Coin: Coin{ID: "c"}, DeveloperScore: 40, PublicInterestScore: 10, MarketCapRank: 4, MoonshotScore: 490},
				{Coin: Coin{ID: "a"}, DeveloperScore: 10, PublicInterestScore: 5, MarketCapRank: 3, MoonshotScore: 325},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := rank(test.cds)
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("got : %v\nwant: %v\n", got, test.expected)
			}
		})
	}
}

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

package main

import (
	"github.com/scharissis/crypto-gem-finder/stonks"
)

func main() {
	s := stonks.NewStonker()
	s.GetGems()
}

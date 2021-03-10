package main

import (
	"log"
	"os"

	"github.com/scharissis/crypto-gem-finder/stonks"
)

func main() {
	s := stonks.NewStonker()
	generateHTML(s, "./web/index.html")

	// For local testing
	// fmt.Printf("serving website from '/web' @ localhost:3000...")
	// log.Fatal(http.ListenAndServe(":3000", http.FileServer(http.Dir("./web"))))
}

func generateHTML(s *stonks.Stonker, path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating file @ %s: %s", path, err)
		return err
	}
	defer f.Close()
	return s.ToHTML(f)
}

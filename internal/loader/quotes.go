package loader

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
)

type Quote struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Length int    `json:"length"`
	Id     int    `json:"id"`
}

type QuoteList struct {
	Language string  `json:"language"`
	Quotes   []Quote `json:"quotes"`
}

//go:embed data/quotes/*
var quotesContent embed.FS

func loadQuote(source string) (Quote, error) {
	filename := fmt.Sprintf("data/quotes/%s.json", source)

	byteValue, err := quotesContent.ReadFile(filename)
	if err != nil {
		return Quote{}, fmt.Errorf("source not found: %s", source)
	}

	var quoteList QuoteList
	if err := json.Unmarshal(byteValue, &quoteList); err != nil {
		return Quote{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(quoteList.Quotes) == 0 {
		return Quote{}, fmt.Errorf("quote file is empty")
	}

	return quoteList.Quotes[rand.Intn(len(quoteList.Quotes))], nil
}

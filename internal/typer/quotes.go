package typer

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

func loadQuote(language string) (Quote, error) {
	filename := fmt.Sprintf("data/quotes/%s.json", language)

	byteValue, err := quotesContent.ReadFile(filename)
	if err != nil {
		return Quote{}, fmt.Errorf("language not found: %s", language)
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

package loader

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/eduqr/ktype/internal/config"
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
	byteValue, err := readQuoteSource(source)
	if err != nil {
		return Quote{}, err
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

func readQuoteSource(source string) ([]byte, error) {
	ktypeDir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	customPath := filepath.Join(ktypeDir, "quotes", source+".json")
	byteValue, err := os.ReadFile(customPath)
	if err == nil {
		return byteValue, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read custom quotes source %q: %w", source, err)
	}

	filename := fmt.Sprintf("data/quotes/%s.json", source)
	byteValue, err = quotesContent.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("quotes source not found: %s", source)
	}

	return byteValue, nil
}

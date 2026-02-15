package typer

import (
	"strings"
)

func GenerateChallenge(language string, wordCount int, quoteMode bool) ([]string, error) {
	if quoteMode {
		quote, err := loadQuote(language)
		if err != nil {
			return nil, err
		}
		return strings.Fields(quote.Text), nil
	}

	return loadWords(language, wordCount)
}

package loader

import (
	"strings"

	"github.com/eduqr/ktype/internal/config"
)

type Challenge struct {
	Words       []string
	QuoteSource string
}

func GenerateChallenge(challenge config.ChallengeConfig) (Challenge, error) {
	if challenge.Quote {
		quote, err := loadQuote(challenge.Source)
		if err != nil {
			return Challenge{}, err
		}
		return Challenge{
			Words:       strings.Fields(quote.Text),
			QuoteSource: strings.Join(strings.Fields(quote.Source), " "),
		}, nil
	}

	words, err := loadWords(challenge.Source, challenge.Words)
	if err != nil {
		return Challenge{}, err
	}
	return Challenge{Words: words}, nil
}

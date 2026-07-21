package loader

import (
	"strings"

	"github.com/eduqr/ktype/internal/config"
)

func GenerateChallenge(challenge config.ChallengeConfig) ([]string, error) {
	if challenge.Quote {
		quote, err := loadQuote(challenge.Source)
		if err != nil {
			return nil, err
		}
		return strings.Fields(quote.Text), nil
	}

	return loadWords(challenge.Source, challenge.Words)
}

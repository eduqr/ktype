package main

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/loader"
	"github.com/eduqr/ktype/internal/logger"
)

func main() {
	config.Parse()
	log := logger.Setup(*config.LogsEnable)

	log.Info("Starting ktype", "wordCount", *config.WordCount, "language", *config.Language, "quoteMode", *config.QuoteMode)

	challenge, err := loader.GenerateChallenge(*config.Language, *config.WordCount, *config.QuoteMode)
	if err != nil {
		log.Error("Failed to generate challenge", "error", err)
		return
	}

	log.Info("Generated challenge", "count", len(challenge), "challenge", challenge)
}

package main

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/debug"
	"github.com/eduqr/ktype/internal/typer"
)

func main() {
	config.Parse()
	logger := debug.Setup(*config.LogsEnable)

	logger.Info("Starting ktype", "wordCount", *config.WordCount, "language", *config.Language, "quoteMode", *config.QuoteMode)

	challenge, err := typer.GenerateChallenge(*config.Language, *config.WordCount, *config.QuoteMode)
	if err != nil {
		logger.Error("Failed to generate challenge", "error", err)
		return
	}

	logger.Info("Generated challenge", "count", len(challenge), "challenge", challenge)
}

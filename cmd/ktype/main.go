package main

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/debug"
)

func main() {
	config.Parse()
	logger := debug.Setup(*config.LogsEnable)

	logger.Info("Starting ktype", "wordCount", *config.WordCount, "language", *config.Language, "phraseMode", *config.PhraseMode)
}
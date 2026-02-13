package main

import (
	"flag"
	"log/slog"
	"os"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonHandler)
	
	logger.Info("ktype cli logs")

	wordCountPtr := flag.Int("w", 25, "The number of words for the challenge")
	languagePtr := flag.String("l", "en", "Challenge's language (options: en, es)")
	phraseModePtr := flag.Bool("p", false, "Use a phrase instead of random words")

	flag.Parse()

	logger.Info("Word count", "value", *wordCountPtr)
	logger.Info("Language", "value", *languagePtr)
	logger.Info("Phrase mode", "value", *phraseModePtr)
}
package config

import "flag"

var (
	WordCount  = flag.Int("w", 25, "The number of words for the challenge")
	Language   = flag.String("l", "en", "Challenge's language (options: en, es)")
	PhraseMode = flag.Bool("p", false, "Use a phrase instead of random words")
	LogsEnable = flag.Bool("debug", false, "Enable logs")
)

func Parse() {
	flag.Parse()
}
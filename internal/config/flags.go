package config

import "flag"

var (
	WordCount  = flag.Int("w", 25, "The number of words for the challenge")
	Language   = flag.String("l", "en", "Challenge's language (options: en, es)")
	QuoteMode = flag.Bool("q", false, "Use a quote instead of random words")
	LogsEnable = flag.Bool("debug", false, "Enable logs")
)

func Parse() {
	flag.Parse()
}
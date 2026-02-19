package config

import "flag"

var (
	WordCount  = flag.Int("w", 25, "The number of words for the challenge")
	Language   = flag.String("l", "english", "Challenge's language (options: english, spanish)")
	QuoteMode  = flag.Bool("q", false, "Use a quote instead of random words")
	LogsEnable = flag.Bool("logs", false, "Enable logs")
)

func Parse() {
	flag.Parse()
}

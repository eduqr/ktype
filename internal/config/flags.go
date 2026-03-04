package config

import (
	"github.com/spf13/cobra"
)

var (
	WordCount  int
	Language   string
	QuoteMode  bool
	LogsEnable bool
)

var RootCmd = &cobra.Command{
	Use:   "ktype",
	Short: "A CLI typing test",
}

func init() {
	RootCmd.Flags().IntVarP(&WordCount, "words", "w", 25, "The number of words for the challenge")
	RootCmd.Flags().StringVarP(&Language, "language", "l", "english", "Challenge's language (options: english, spanish)")
	RootCmd.Flags().BoolVarP(&QuoteMode, "quote", "q", false, "Use a quote instead of random words")
	RootCmd.Flags().BoolVar(&LogsEnable, "logs", false, "Enable logs")
}

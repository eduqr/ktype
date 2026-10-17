package cli

import "github.com/spf13/viper"

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help for ktype")
	rootCmd.Flags().IntP("words", "w", 25, "The number of words for the challenge")
	rootCmd.Flags().StringP("source", "s", "english", "Challenge source")
	rootCmd.Flags().BoolP("quote", "q", false, "Use a quote instead of random words")

	_ = viper.BindPFlag("challenge.words", rootCmd.Flags().Lookup("words"))
	_ = viper.BindPFlag("challenge.source", rootCmd.Flags().Lookup("source"))
	_ = viper.BindPFlag("challenge.quote", rootCmd.Flags().Lookup("quote"))
}

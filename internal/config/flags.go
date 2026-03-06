package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "ktype",
	Short: "A CLI typing test",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitConfig()
	},
}

func init() {
	RootCmd.Flags().IntP("words", "w", 25, "The number of words for the challenge")
	RootCmd.Flags().StringP("language", "l", "english", "Challenge's language (options: english, spanish)")
	RootCmd.Flags().BoolP("quote", "q", false, "Use a quote instead of random words")
	RootCmd.Flags().Bool("logs", false, "Enable logs")

	_ = viper.BindPFlag("challenge.words", RootCmd.Flags().Lookup("words"))
	_ = viper.BindPFlag("challenge.language", RootCmd.Flags().Lookup("language"))
	_ = viper.BindPFlag("challenge.quote", RootCmd.Flags().Lookup("quote"))
	_ = viper.BindPFlag("logs.enabled", RootCmd.Flags().Lookup("logs"))
}

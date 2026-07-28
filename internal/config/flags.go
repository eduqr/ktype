package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:           "ktype",
	Short:         "A CLI typing test",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitConfig()
	},
}

func init() {
	RootCmd.Flags().IntP("words", "w", 25, "The number of words for the challenge")
	RootCmd.Flags().StringP("source", "s", "english", "Challenge source")
	RootCmd.Flags().BoolP("quote", "q", false, "Use a quote instead of random words")
	RootCmd.Flags().Bool("logs", false, "Enable diagnostic logs")

	_ = viper.BindPFlag("challenge.words", RootCmd.Flags().Lookup("words"))
	_ = viper.BindPFlag("challenge.source", RootCmd.Flags().Lookup("source"))
	_ = viper.BindPFlag("challenge.quote", RootCmd.Flags().Lookup("quote"))
	_ = viper.BindPFlag("logs.enabled", RootCmd.Flags().Lookup("logs"))
}

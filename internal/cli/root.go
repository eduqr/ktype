package cli

import (
	"github.com/eduqr/ktype/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "ktype",
	Short:         "Typing test for the terminal",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}
		return config.InitConfig()
	},
}

package cli

import (
	"fmt"
	"os"

	"github.com/eduqr/ktype/internal/config"
	"github.com/eduqr/ktype/internal/session"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.RunE = run
	rootCmd.AddCommand(listCmd, versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	cfg := config.Cfg

	session, err := session.New(cfg)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run()
}

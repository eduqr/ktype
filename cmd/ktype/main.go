package main

import (
	"fmt"
	"os"

	"github.com/eduqr/ktype/internal/config"
)

func main() {
	config.RootCmd.RunE = run
	config.RootCmd.AddCommand(listCmd)
	config.RootCmd.AddCommand(versionCmd)

	if err := config.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

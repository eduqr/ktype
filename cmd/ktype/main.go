package main

import (
	"fmt"
	"os"

	"github.com/eduqr/ktype/internal/config"
)

func main() {
	config.RootCmd.RunE = run
	config.RootCmd.AddCommand(listCmd)

	if err := config.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

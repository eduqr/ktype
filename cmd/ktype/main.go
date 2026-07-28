package main

import (
	"fmt"
	"os"

	"github.com/eduqr/ktype/internal/config"
)

func main() {
	config.RootCmd.Run = run
	config.RootCmd.AddCommand(listCmd)

	if err := config.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

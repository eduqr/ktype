package main

import (
	"github.com/eduqr/ktype/internal/loader"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ktype sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		loader.ListSources()
		return nil
	},
}

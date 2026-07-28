package main

import (
	"fmt"
	"io"

	"github.com/eduqr/ktype/internal/loader"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ktype sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		sources, err := loader.ListSources()
		if err != nil {
			return err
		}

		printSources(cmd, sources)
		return nil
	},
}

func printSources(cmd *cobra.Command, sources loader.Sources) {
	width := maxNameLen(sources)
	output := cmd.OutOrStdout()

	fmt.Fprintln(output, "Words:")
	for _, source := range sources.Words {
		printSource(output, width, source)
	}

	fmt.Fprintln(output, "\nQuotes:")
	for _, source := range sources.Quotes {
		printSource(output, width, source)
	}
}

func printSource(output io.Writer, width int, source loader.Source) {
	origin := source.Origin
	if source.Overrides {
		origin += " (overrides built-in)"
	}

	fmt.Fprintf(output, "  %-*s  %s\n", width, source.Name, origin)
}

func maxNameLen(sources loader.Sources) int {
	longest := 0

	for _, source := range append(sources.Words, sources.Quotes...) {
		longest = max(longest, len(source.Name))
	}

	return longest + 2
}

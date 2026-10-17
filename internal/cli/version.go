package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the ktype version",
	Run: func(cmd *cobra.Command, args []string) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "ktype (devel)")
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "ktype %s\n", info.Main.Version)
	},
}

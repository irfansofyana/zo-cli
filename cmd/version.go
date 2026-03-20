package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is the current CLI version.
// Priority (highest to lowest):
//  1. Set via -ldflags at build time (e.g. goreleaser / custom Makefile)
//  2. Embedded by the Go toolchain when installed via `go install module@version`
//  3. Falls back to "dev" for local builds with no version info
var Version = "dev"

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
	rootCmd.Version = Version
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of zo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), Version)
	},
}

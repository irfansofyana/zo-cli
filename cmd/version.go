package cmd

import "runtime/debug"

// Version is the current CLI version.
// Set via -ldflags at build time, or read from build info when
// installed via `go install module@version`. Falls back to "dev".
var Version = "dev"

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
	rootCmd.Version = Version
}

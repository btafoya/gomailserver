package main

import (
	"fmt"
	"os"

	"github.com/btafoya/gomailserver/internal/commands"
)

// Injected at build time via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Build version string with commit and date info when available
	versionInfo := version
	if commit != "none" && len(commit) >= 7 {
		versionInfo = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit[:7], date)
	}

	if err := commands.Execute(versionInfo); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/Kocoro-lab/ptengine-cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	if err := cmd.Execute(); err != nil {
		if exitErr, ok := err.(*cmd.ExitError); ok {
			os.Exit(exitErr.Code)
		}
		os.Exit(1)
	}
}

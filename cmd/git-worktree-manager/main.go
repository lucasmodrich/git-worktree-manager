package main

import (
	"os"

	"github.com/lucasmodrich/git-worktree-manager/internal/commands"
)

// Build info injected via ldflags:
//
//	-X main.version=x.y.z
//	-X main.commit=abc1234
//	-X main.date=2006-01-02T15:04:05Z
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.SetBuildInfo(version, commit, date)

	// Execute root command
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

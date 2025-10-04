package main

import (
	"os"

	"github.com/lucasmodrich/git-worktree-manager/internal/commands"
)

// version will be set via ldflags during build: -X main.version=x.y.z
var version = "dev"

func main() {
	// Set version in root command
	commands.SetVersion(version)

	// Execute root command
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

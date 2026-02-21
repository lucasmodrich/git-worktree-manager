package ui

import (
	"fmt"
)

// PrintStatus prints a status message with an emoji prefix
func PrintStatus(emoji, message string) {
	fmt.Printf("%s %s\n", emoji, message)
}

// PrintDryRun prints a dry-run prefixed message
func PrintDryRun(message string) {
	fmt.Printf("🔍 [DRY-RUN] %s\n", message)
}

// PrintProgress prints a progress message (for long-running operations)
func PrintProgress(message string) {
	fmt.Println(message)
}

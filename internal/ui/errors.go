package ui

import (
	"fmt"
	"os"
)

// FormatError formats an error message with actionable guidance
// Returns a formatted string with ❌ emoji for error and 💡 emoji for guidance
func FormatError(err error, guidance string) string {
	return fmt.Sprintf("❌ %s\n💡 %s", err.Error(), guidance)
}

// PrintError prints a formatted error message to stderr
func PrintError(err error, guidance string) {
	fmt.Fprintln(os.Stderr, FormatError(err, guidance))
}

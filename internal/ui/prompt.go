package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptYesNo prompts the user with a yes/no question
// Returns true if user answers yes, false if no
func PromptYesNo(question string) (bool, error) {
	fmt.Printf("%s [y/N]: ", question)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}
		// EOF reached, default to no
		return false, nil
	}

	answer := strings.TrimSpace(scanner.Text())
	return parseYesNo(answer)
}

// parseYesNo parses a yes/no answer string
func parseYesNo(answer string) (bool, error) {
	answer = strings.ToLower(strings.TrimSpace(answer))

	switch answer {
	case "y", "yes":
		return true, nil
	case "n", "no", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: expected y/n, got %q", answer)
	}
}

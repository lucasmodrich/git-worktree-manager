package ui

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFormatError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		guidance     string
		wantContains []string
	}{
		{
			name:         "error with guidance",
			err:          errors.New("file not found"),
			guidance:     "Check the file path and try again",
			wantContains: []string{"❌", "file not found", "💡", "Check the file path and try again"},
		},
		{
			name:         "bare directory exists error",
			err:          errors.New(".bare directory already exists"),
			guidance:     "Remove existing .bare directory or run setup in a different directory",
			wantContains: []string{"❌", ".bare directory already exists", "💡", "Remove existing .bare directory"},
		},
		{
			name:         "not in worktree repo error",
			err:          errors.New("not in a worktree-managed repository"),
			guidance:     "Run this command from a directory where .git points to .bare",
			wantContains: []string{"❌", "not in a worktree-managed repository", "💡", ".git points to .bare"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatError(tt.err, tt.guidance)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FormatError() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

func TestPrintError(t *testing.T) {
	// Capture stderr via a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	old := os.Stderr
	os.Stderr = w

	PrintError(errors.New("test error"), "test guidance")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	out := buf.String()

	if !strings.Contains(out, "test error") {
		t.Errorf("PrintError() stderr = %q, want to contain %q", out, "test error")
	}
	if !strings.Contains(out, "test guidance") {
		t.Errorf("PrintError() stderr = %q, want to contain %q", out, "test guidance")
	}
}

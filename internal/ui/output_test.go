package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintStatus(t *testing.T) {
	tests := []struct {
		name          string
		emoji         string
		message       string
		wantContains  []string
	}{
		{
			name:         "success status",
			emoji:        "✅",
			message:      "Setup complete!",
			wantContains: []string{"✅", "Setup complete!"},
		},
		{
			name:         "error status",
			emoji:        "❌",
			message:      "Operation failed",
			wantContains: []string{"❌", "Operation failed"},
		},
		{
			name:         "info status",
			emoji:        "📡",
			message:      "Fetching from origin",
			wantContains: []string{"📡", "Fetching from origin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintStatus(tt.emoji, tt.message)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("PrintStatus() output = %q, want to contain %q", output, want)
				}
			}
		})
	}
}

func TestPrintDryRun(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		wantContains []string
	}{
		{
			name:         "dry-run message",
			message:      "Would create branch",
			wantContains: []string{"🔍", "[DRY-RUN]", "Would create branch"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintDryRun(tt.message)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("PrintDryRun() output = %q, want to contain %q", output, want)
				}
			}
		})
	}
}

func TestPrintProgress(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "progress message",
			message: "Cloning repository...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintProgress(tt.message)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tt.message) {
				t.Errorf("PrintProgress() output = %q, want to contain %q", output, tt.message)
			}
		})
	}
}

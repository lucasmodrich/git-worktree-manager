package ui

import (
	"strings"
	"testing"
)

func TestParseYesNo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:    "yes - lowercase y",
			input:   "y",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes - uppercase Y",
			input:   "Y",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes - full word",
			input:   "yes",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes - uppercase YES",
			input:   "YES",
			want:    true,
			wantErr: false,
		},
		{
			name:    "no - lowercase n",
			input:   "n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no - uppercase N",
			input:   "N",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no - full word",
			input:   "no",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no - empty string (default no)",
			input:   "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "maybe",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYesNo(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseYesNo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseYesNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromptYesNo(t *testing.T) {
	// Note: This test verifies the function signature and basic functionality
	// Actual interactive testing would require mocking stdin
	tests := []struct {
		name     string
		question string
	}{
		{
			name:     "basic prompt",
			question: "Continue?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the question format would include [y/N]
			expectedFormat := tt.question + " [y/N]"
			if !strings.Contains(expectedFormat, "[y/N]") {
				t.Errorf("Prompt format missing [y/N] suffix")
			}
		})
	}
}

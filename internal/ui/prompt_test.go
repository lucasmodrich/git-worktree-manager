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
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:  "yes answer",
			input: "y\n",
			want:  true,
		},
		{
			name:  "no answer",
			input: "n\n",
			want:  false,
		},
		{
			name:  "empty input defaults to no",
			input: "\n",
			want:  false,
		},
		{
			name:  "EOF defaults to no",
			input: "",
			want:  false,
		},
		{
			name:    "invalid input returns error",
			input:   "maybe\n",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PromptYesNo("Continue?", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("PromptYesNo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PromptYesNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

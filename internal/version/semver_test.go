package version

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "simple version",
			input: "1.0.0",
			want: &Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: nil,
				Build:      "",
				Original:   "1.0.0",
			},
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want: &Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: nil,
				Build:      "",
				Original:   "v1.2.3",
			},
		},
		{
			name:  "version with prerelease",
			input: "1.0.0-alpha",
			want: &Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: []string{"alpha"},
				Build:      "",
				Original:   "1.0.0-alpha",
			},
		},
		{
			name:  "version with numeric prerelease",
			input: "1.0.0-alpha.1",
			want: &Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: []string{"alpha", "1"},
				Build:      "",
				Original:   "1.0.0-alpha.1",
			},
		},
		{
			name:  "version with build metadata",
			input: "1.0.0+build.1",
			want: &Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: nil,
				Build:      "build.1",
				Original:   "1.0.0+build.1",
			},
		},
		{
			name:  "version with prerelease and build",
			input: "1.0.0-rc.1+build",
			want: &Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: []string{"rc", "1"},
				Build:      "build",
				Original:   "1.0.0-rc.1+build",
			},
		},
		{
			name:    "invalid version - missing patch",
			input:   "1.0",
			wantErr: true,
		},
		{
			name:    "invalid version - non-numeric",
			input:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch {
					t.Errorf("ParseVersion() version = %d.%d.%d, want %d.%d.%d",
						got.Major, got.Minor, got.Patch, tt.want.Major, tt.want.Minor, tt.want.Patch)
				}
				if got.Build != tt.want.Build {
					t.Errorf("ParseVersion() build = %v, want %v", got.Build, tt.want.Build)
				}
			}
		})
	}
}

func TestVersionGreaterThan(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want bool
	}{
		// Equal versions
		{name: "1.0.0 == 1.0.0", v1: "1.0.0", v2: "1.0.0", want: false},
		{name: "1.1.5 == 1.1.5", v1: "1.1.5", v2: "1.1.5", want: false},

		// Patch version comparison
		{name: "1.0.1 > 1.0.0", v1: "1.0.1", v2: "1.0.0", want: true},
		{name: "1.1.6 > 1.1.5", v1: "1.1.6", v2: "1.1.5", want: true},

		// Minor version comparison
		{name: "1.1.0 > 1.0.99", v1: "1.1.0", v2: "1.0.99", want: true},

		// Major version comparison
		{name: "2.0.0 < 10.0.0", v1: "2.0.0", v2: "10.0.0", want: false},

		// v prefix handling (should be treated as equal)
		{name: "v1.2.3 == 1.2.3", v1: "v1.2.3", v2: "1.2.3", want: false},

		// Prerelease precedence (release > prerelease)
		{name: "1.0.0-alpha < 1.0.0", v1: "1.0.0-alpha", v2: "1.0.0", want: false},
		{name: "1.0.0 > 1.0.0-alpha", v1: "1.0.0", v2: "1.0.0-alpha", want: true},
		{name: "1.0.0 > 1.0.0-rc.1", v1: "1.0.0", v2: "1.0.0-rc.1", want: true},

		// Prerelease comparison
		{name: "1.0.0-alpha.1 > 1.0.0-alpha", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha", want: true},

		// Alphanumeric vs numeric prerelease (alphanumeric > numeric in ASCII order)
		{name: "1.0.0-alpha.beta > 1.0.0-alpha.1", v1: "1.0.0-alpha.beta", v2: "1.0.0-alpha.1", want: true},

		// Numeric prerelease comparison
		{name: "1.0.0-beta.11 > 1.0.0-beta.2", v1: "1.0.0-beta.11", v2: "1.0.0-beta.2", want: true},

		// Alphanumeric prerelease > numeric prerelease
		{name: "1.0.0-alpha > 1.0.0-1", v1: "1.0.0-alpha", v2: "1.0.0-1", want: true},

		// Build metadata ignored
		{name: "1.0.0+build.1 == 1.0.0+other (build ignored)", v1: "1.0.0+build.1", v2: "1.0.0+other", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver1, err1 := ParseVersion(tt.v1)
			ver2, err2 := ParseVersion(tt.v2)

			if err1 != nil || err2 != nil {
				t.Fatalf("ParseVersion failed: v1=%v, v2=%v", err1, err2)
			}

			got := ver1.GreaterThan(ver2)
			if got != tt.want {
				t.Errorf("GreaterThan() = %v, want %v (comparing %s vs %s)", got, tt.want, tt.v1, tt.v2)
			}
		})
	}
}

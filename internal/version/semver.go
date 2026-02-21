package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version (semver 2.0.0 compliant)
type Version struct {
	Major      int      // Major version number
	Minor      int      // Minor version number
	Patch      int      // Patch version number
	Prerelease []string // Prerelease identifiers (split on '.')
	Build      string   // Build metadata (ignored in comparison)
	Original   string   // Original version string
}

var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-\.]+))?(?:\+([0-9A-Za-z\-\.]+))?$`)

// ParseVersion parses a semantic version string
func ParseVersion(v string) (*Version, error) {
	matches := semverRegex.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version: %s", v)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	var prerelease []string
	if matches[4] != "" {
		prerelease = strings.Split(matches[4], ".")
	}

	build := matches[5]

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
		Build:      build,
		Original:   v,
	}, nil
}

// GreaterThan returns true if this version is greater than the other version
// Implements semver 2.0.0 comparison rules, matching the Bash version_gt() behavior
func (v *Version) GreaterThan(other *Version) bool {
	// Compare major, minor, patch
	if v.Major > other.Major {
		return true
	}
	if v.Major < other.Major {
		return false
	}

	if v.Minor > other.Minor {
		return true
	}
	if v.Minor < other.Minor {
		return false
	}

	if v.Patch > other.Patch {
		return true
	}
	if v.Patch < other.Patch {
		return false
	}

	// Main version parts are equal, handle prerelease precedence
	// Release (no prerelease) > Prerelease
	hasPreV := len(v.Prerelease) > 0
	hasPreOther := len(other.Prerelease) > 0

	if !hasPreV && !hasPreOther {
		return false // Equal
	}
	if !hasPreV && hasPreOther {
		return true // Release > prerelease
	}
	if hasPreV && !hasPreOther {
		return false // Prerelease < release
	}

	// Both have prerelease, compare identifiers
	return comparePrereleaseIdentifiers(v.Prerelease, other.Prerelease)
}

// comparePrereleaseIdentifiers compares two prerelease identifier arrays
func comparePrereleaseIdentifiers(a, b []string) bool {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i < maxLen; i++ {
		// If one array is shorter, shorter < longer
		if i >= len(a) {
			return false // a is shorter, so a < b
		}
		if i >= len(b) {
			return true // b is shorter, so a > b
		}

		idA := a[i]
		idB := b[i]

		numA, errA := strconv.Atoi(idA)
		numB, errB := strconv.Atoi(idB)

		// Both numeric: compare numerically
		if errA == nil && errB == nil {
			if numA > numB {
				return true
			}
			if numA < numB {
				return false
			}
			// Equal, continue to next identifier
			continue
		}

		// Numeric identifiers have lower precedence than non-numeric
		if errA == nil && errB != nil {
			return false // numeric < alphanumeric
		}
		if errA != nil && errB == nil {
			return true // alphanumeric > numeric
		}

		// Both non-numeric: ASCII lexical comparison
		if idA > idB {
			return true
		}
		if idA < idB {
			return false
		}
		// Equal, continue to next identifier
	}

	// All identifiers equal
	return false
}

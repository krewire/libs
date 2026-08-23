package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version is a semantic version per https://semver.org/.
// Build metadata is retained but ignored for precedence.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

var versionRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?(?:\+([0-9A-Za-z.-]+))?$`)

// ParseVersion parses s as a semantic version. Leading "v" is optional.
func ParseVersion(s string) (Version, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Version{}, UsageError("version is required")
	}
	m := versionRe.FindStringSubmatch(s)
	if m == nil {
		return Version{}, UsageError(fmt.Sprintf("invalid version %q: want MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]", s))
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	return Version{Major: major, Minor: minor, Patch: patch, PreRelease: m[4], Build: m[5]}, nil
}

// MustParseVersion parses s or panics. Use for constants.
func MustParseVersion(s string) Version {
	v, err := ParseVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns the canonical string form without leading "v".
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		s += "-" + v.PreRelease
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}

// Compare returns -1 if v < other, 0 if equal, 1 if v > other per semver precedence.
// Build metadata is ignored.
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	// Pre-release: absence has higher precedence than presence
	if v.PreRelease == "" && other.PreRelease != "" {
		return 1
	}
	if v.PreRelease != "" && other.PreRelease == "" {
		return -1
	}
	if v.PreRelease != other.PreRelease {
		if v.PreRelease < other.PreRelease {
			return -1
		}
		return 1
	}
	return 0
}

// Less reports whether v < other.
func (v Version) Less(other Version) bool { return v.Compare(other) < 0 }

// Equal reports whether v == other (ignoring build).
func (v Version) Equal(other Version) bool { return v.Compare(other) == 0 }

// IsCompatible reports whether actual satisfies required per semver caret semantics
// for the Krewire ecosystem: for 0.y.z, minor must match; for >=1.0.0, major must match and actual >= required.
func (v Version) IsCompatible(required Version) bool {
	if required.Major == 0 {
		// 0.y.z: minor must match, patch >= required
		if v.Major != 0 || v.Minor != required.Minor {
			return false
		}
		return v.Compare(required) >= 0
	}
	if v.Major != required.Major {
		return false
	}
	return v.Compare(required) >= 0
}

// ModuleName identifies a Krewire module in the ecosystem compatibility matrix.
type ModuleName string

const (
	ModuleFramework ModuleName = "framework"
	ModuleLibs      ModuleName = "libs"
	ModuleMdbind    ModuleName = "mdbind"
	ModuleKrewire   ModuleName = "krewire"
	ModuleGuild     ModuleName = "guild"
	ModuleDocs      ModuleName = "docs"
	ModuleLanding   ModuleName = "krewire.github.io"
	ModuleInternal  ModuleName = "internal"
)

// CurrentVersion is the libs module's own version. Bump per release.
var CurrentVersion = MustParseVersion("0.1.0")

// EcosystemVersions is the known-good compatibility matrix for the current release.
var EcosystemVersions = map[ModuleName]Version{
	ModuleFramework: MustParseVersion("0.1.0"),
	ModuleLibs:      CurrentVersion,
	ModuleMdbind:    MustParseVersion("0.1.0"),
	ModuleKrewire:   MustParseVersion("0.1.0"),
	ModuleGuild:     MustParseVersion("0.1.0"),
	ModuleInternal:  MustParseVersion("0.1.0"),
}

// CheckEcosystemCompatibility verifies that actual versions satisfy required versions per IsCompatible.
// Pass the go.mod require versions as actual; the matrix as required.
func CheckEcosystemCompatibility(required, actual map[ModuleName]Version) error {
	for mod, reqVer := range required {
		actVer, ok := actual[mod]
		if !ok {
			return FailureError(fmt.Sprintf("missing version for module %q", mod))
		}
		if !actVer.IsCompatible(reqVer) {
			return FailureError(fmt.Sprintf("module %q version %s is not compatible with required %s", mod, actVer.String(), reqVer.String()))
		}
	}
	return nil
}

// Versioned is implemented by any module that exposes its version via core.Version.
type Versioned interface {
	Version() Version
}

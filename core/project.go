package core

import (
	"fmt"
	"regexp"
	"strings"
)

// Project describes a Krewire project for validation.
type Project struct {
	Name       string `json:"name"`
	ModulePath string `json:"modulePath"`
	Kind       Kind   `json:"kind"`
	ConfigPath string `json:"configPath"`
}

var nameRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var moduleRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._/\-]*$`)

// Validate checks project invariants.
func (p Project) Validate() error {
	if p.Name == "" {
		return UsageError("project name is required")
	}
	if !nameRe.MatchString(p.Name) {
		return UsageError(fmt.Sprintf("invalid project name %q: want kebab-case ^[a-z][a-z0-9-]*$", p.Name))
	}
	if p.ModulePath != "" && !moduleRe.MatchString(p.ModulePath) {
		return UsageError(fmt.Sprintf("invalid module path %q", p.ModulePath))
	}
	if !p.Kind.IsValid() {
		return UsageError(fmt.Sprintf("invalid project kind %q", p.Kind))
	}
	if p.ConfigPath != "" {
		if err := ValidateKrewireYamlPath(p.ConfigPath); err != nil {
			return err
		}
	}
	return nil
}

// ValidateKrewireYamlPath ensures the config path is krewire.yaml.
func ValidateKrewireYamlPath(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return UsageError("config path is required")
	}
	// Normalize: allow ./krewire.yaml, krewire.yaml, /abs/krewire.yaml
	base := path
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	if base != "krewire.yaml" {
		return UsageError(fmt.Sprintf("invalid config path %q: must be krewire.yaml (no ssg.yaml)", path))
	}
	return nil
}

// optInBlockedPackages are the framework batteries a KindApp monolith must not
// import directly: they are opt-in batteries reserved for their own project
// kinds. Import paths are matched on a path-segment boundary so subpackages
// (e.g. "framework/service/gateway") are covered without matching unrelated
// packages that merely share a prefix.
var optInBlockedPackages = []string{
	"github.com/krewire/framework/service",
	"github.com/krewire/framework/infra",
}

// HasOptInViolation reports whether importing the given import paths violates
// opt-in for the declared kind. Per KWL-CORE-022, a KindApp monolith importing
// framework/service or framework/infra is a violation: those batteries belong to
// the service and infra kinds and must cost a monolith nothing.
//
// framework/worker is deliberately not blocked — a KindApp project may run
// background jobs in-process — and the batteries whose Kind matches the project
// are never violations (e.g. KindService importing framework/service).
func HasOptInViolation(kind Kind, imported []string) bool {
	if kind != KindApp {
		return false
	}
	for _, imp := range imported {
		for _, blocked := range optInBlockedPackages {
			if imp == blocked || strings.HasPrefix(imp, blocked+"/") {
				return true
			}
		}
	}
	return false
}

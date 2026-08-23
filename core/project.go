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

// IsOptIn reports whether importing the given import paths violates opt-in
// for the declared kind. For example, a KindApp monolith importing
// framework/service should be flagged.
func IsOptIn(kind Kind, imported []string) bool {
	// Only app is expected to be zero-cost; worker/service/infra are opt-in by kind.
	if kind != KindApp {
		return true
	}
	for _, imp := range imported {
		if strings.Contains(imp, "framework/service") || strings.Contains(imp, "framework/infra") || strings.Contains(imp, "framework/worker") && !strings.Contains(imp, "runtime") {
			// runtime is allowed for app (frontend), but service/infra/worker are opt-in
			if strings.Contains(imp, "framework/service") || strings.Contains(imp, "framework/infra") {
				return false
			}
		}
	}
	return true
}

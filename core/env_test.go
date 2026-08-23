// Tests for KWL-K4T7W
package core

import (
	"errors"
	"strings"
	"testing"
)

// Spec: KWL-K4T7W KWL-ENVV-001 Scope: Package
func TestKWL_ENVV_001_ParseEnv_KnownValuesAndDefault(t *testing.T) {
	cases := map[string]Env{
		"":             EnvLocal,
		"   ":          EnvLocal,
		"local":        EnvLocal,
		"production":   EnvProduction,
		"testing":      EnvTesting,
		"LOCAL":        EnvLocal,
		" Production ": EnvProduction,
	}
	for in, want := range cases {
		got, err := ParseEnv(in)
		if err != nil || got != want {
			t.Errorf("ParseEnv(%q) = (%q, %v), want (%q, nil)", in, got, err, want)
		}
	}
	if envs := Envs(); len(envs) != 3 {
		t.Errorf("Envs() = %v, want exactly three values", envs)
	}
}

// Spec: KWL-K4T7W KWL-ENVV-002 Scope: Package
func TestKWL_ENVV_002_ParseEnv_RejectsUnknownWithUsageError(t *testing.T) {
	_, err := ParseEnv("staging")
	if err == nil {
		t.Fatal("ParseEnv(staging) = nil error, want usage error")
	}
	var ce interface{ ExitCode() ExitCode }
	if !errors.As(err, &ce) || ce.ExitCode() != ExitCodeUsage {
		t.Errorf("ParseEnv(staging) error = %v, want usage exit code", err)
	}
	for _, allowed := range []string{"local", "production", "testing"} {
		if !strings.Contains(err.Error(), allowed) {
			t.Errorf("ParseEnv error %q should name allowed value %q", err.Error(), allowed)
		}
	}
}

package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmpty(t *testing.T) {
	c, err := Load("")
	if err != nil {
		t.Fatalf("Load empty path error: %v", err)
	}
	if len(c) != 0 {
		t.Errorf("expected empty config, got %v", c)
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpdir := t.TempDir()
	path := filepath.Join(tmpdir, "nonexistent.yaml")
	c, err := Load(path)
	if err != nil {
		t.Fatalf("Load non-existent error: %v", err)
	}
	if len(c) != 0 {
		t.Errorf("expected empty config for non-existent file, got %v", c)
	}
}

func TestLoadAndSave(t *testing.T) {
	tmpdir := t.TempDir()
	path := filepath.Join(tmpdir, "config.yaml")

	// Create test config
	original := Config{
		"app.name":       "testapp",
		"server.port":    "8080",
		"database.host":  "localhost",
		"features.debug": "true",
	}

	if err := original.Save(path); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	expected := map[string]string{
		"app.name":       "testapp",
		"server.port":    "8080",
		"database.host":  "localhost",
		"features.debug": "true",
	}
	for k, v := range expected {
		if loaded[k] != v {
			t.Errorf("key %q: got %q, want %q", k, loaded[k], v)
		}
	}
}

func TestGetSetDelete(t *testing.T) {
	c := make(Config)
	c.Set("key1", "value1")
	if c.Get("key1") != "value1" {
		t.Errorf("Get: got %q, want value1", c.Get("key1"))
	}
	c.Delete("key1")
	if c.Get("key1") != "" {
		t.Errorf("Delete failed: got %q, want empty", c.Get("key1"))
	}
}

func TestKeysSorted(t *testing.T) {
	c := Config{
		"zebra": "1",
		"alpha": "2",
		"beta":  "3",
	}
	keys := c.Keys()
	expected := []string{"alpha", "beta", "zebra"}
	if len(keys) != 3 || keys[0] != "alpha" || keys[1] != "beta" || keys[2] != "zebra" {
		t.Errorf("Keys not sorted: got %v, want %v", keys, expected)
	}
}

func TestMerge(t *testing.T) {
	a := Config{"a": "1", "b": "2"}
	b := Config{"b": "new", "c": "3"}
	merged := a.Merge(b)
	expected := map[string]string{"a": "1", "b": "new", "c": "3"}
	for k, v := range expected {
		if merged[k] != v {
			t.Errorf("Merge: key %q got %q, want %q", k, merged[k], v)
		}
	}
}

func TestWithDefaults(t *testing.T) {
	c := Config{"existing": "value"}
	defaults := Config{"existing": "default", "new": "default"}
	result := c.WithDefaults(defaults)
	if result["existing"] != "value" {
		t.Errorf("existing key overwritten: got %q, want value", result["existing"])
	}
	if result["new"] != "default" {
		t.Errorf("default not applied: got %q, want default", result["new"])
	}
}

func TestEnvOverride(t *testing.T) {
	os.Setenv("APP_SERVER_PORT", "8080")
	os.Setenv("APP_DB_HOST", "prod-db")
	defer os.Unsetenv("APP_SERVER_PORT")
	defer os.Unsetenv("APP_DB_HOST")

	c := make(Config)
	c.EnvOverride("APP")
	if c.Get("server.port") != "8080" {
		t.Errorf("EnvOverride: server.port = %q, want 8080", c.Get("server.port"))
	}
	if c.Get("db.host") != "prod-db" {
		t.Errorf("EnvOverride: db.host = %q, want prod-db", c.Get("db.host"))
	}
}

func TestFlattenUnflatten(t *testing.T) {
	nested := map[string]any{
		"server": map[string]any{
			"port": 8080,
			"host": "localhost",
		},
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"features": []any{
			map[string]any{"name": "auth", "enabled": true},
			map[string]any{"name": "cache", "enabled": false},
		},
	}

	flat := make(Config)
	flatten(nested, "", flat)

	unflat := unflatten(flat)

	// Values are stored as strings, so we need to handle type conversion
	server := unflat["server"].(map[string]any)
	if server["port"] != "8080" && server["port"] != float64(8080) {
		t.Errorf("unflatten port mismatch: got %v (%T)", server["port"], server["port"])
	}
	if server["host"] != "localhost" {
		t.Errorf("unflatten host mismatch")
	}
}

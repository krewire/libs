package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the key-value configuration store.
// It's a simple string-to-string map that can be loaded from and saved to YAML files.
type Config map[string]string

// Load reads a YAML file and returns a Config map.
// If the file doesn't exist, it returns an empty Config (not an error).
func Load(path string) (Config, error) {
	cfg := make(Config)
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("cfg: read %q: %w", path, err)
	}
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("cfg: parse %q: %w", path, err)
	}
	flatten(raw, "", cfg)
	return cfg, nil
}

// Save writes the Config map to a YAML file.
func (c Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cfg: mkdir %q: %w", dir, err)
	}
	// Convert flat map back to nested structure for readable YAML
	nested := unflatten(c)
	data, err := yaml.Marshal(nested)
	if err != nil {
		return fmt.Errorf("cfg: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Get returns the value for key, or empty string if not found.
func (c Config) Get(key string) string {
	return c[key]
}

// Set sets a key-value pair.
func (c Config) Set(key, value string) {
	c[key] = value
}

// Delete removes a key.
func (c Config) Delete(key string) {
	delete(c, key)
}

// Keys returns all keys in sorted order.
func (c Config) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	// Sort for deterministic output
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

// flatten converts a nested map to a flat map with dot-notation keys.
func flatten(m map[string]any, prefix string, out Config) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch v := v.(type) {
		case map[string]any:
			flatten(v, key, out)
		case []any:
			// For slices, we store as JSON-like string or indexed keys
			for i, item := range v {
				if m, ok := item.(map[string]any); ok {
					flatten(m, fmt.Sprintf("%s.%d", key, i), out)
				} else {
					out[fmt.Sprintf("%s.%d", key, i)] = fmt.Sprintf("%v", item)
				}
			}
		default:
			out[key] = fmt.Sprintf("%v", v)
		}
	}
}

// unflatten converts a flat map with dot-notation keys back to a nested map.
func unflatten(flat Config) map[string]any {
	result := make(map[string]any)
	for k, v := range flat {
		parts := strings.Split(k, ".")
		current := result
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = v
			} else {
				if next, ok := current[part].(map[string]any); ok {
					current = next
				} else {
					next = make(map[string]any)
					current[part] = next
					current = next
				}
			}
		}
	}
	return result
}

// Merge merges another Config into this one, with other taking precedence.
func (c Config) Merge(other Config) Config {
	result := make(Config)
	for k, v := range c {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}

// EnvOverride applies environment variables with the given prefix.
// Environment variables are converted to lower-case keys with dots.
// e.g., APP_SERVER_PORT=8080 -> server.port
func (c Config) EnvOverride(prefix string) Config {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k, value := parts[0], parts[1]
		if !strings.HasPrefix(k, prefix+"_") {
			continue
		}
		// Remove prefix and convert to lower-case with dots
		key := strings.ToLower(strings.TrimPrefix(k, prefix+"_"))
		key = strings.ReplaceAll(key, "_", ".")
		c[key] = value
	}
	return c
}

// Defaults returns a new Config with default values applied (non-destructive).
func (c Config) WithDefaults(defaults Config) Config {
	result := make(Config)
	for k, v := range c {
		result[k] = v
	}
	for k, v := range defaults {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	return result
}

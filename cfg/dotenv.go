package cfg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetOr returns the value for key, or fallback when the key is absent.
// It is the lookup primitive for config getters that read environment
// variables with in-code defaults (KWN-Q7X4M).
func (c Config) GetOr(key, fallback string) string {
	if v, ok := c[key]; ok && v != "" {
		return v
	}
	return fallback
}

// LoadDotEnv parses the .env file at path and exports its KEY=VALUE pairs
// into the process environment (KWL-Q7X4M KWL-DOTV-001). Variables already
// present in the process environment are never overwritten, so real
// environment wins over .env. A missing file is not an error.
func LoadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cfg: read %q: %w", path, err)
	}
	pairs, err := ParseDotEnv(data)
	if err != nil {
		return fmt.Errorf("cfg: parse %q: %w", path, err)
	}
	for _, kv := range pairs {
		if _, exists := os.LookupEnv(kv.key); !exists {
			if err := os.Setenv(kv.key, kv.value); err != nil {
				return fmt.Errorf("cfg: setenv %q: %w", kv.key, err)
			}
		}
	}
	return nil
}

type dotenvPair struct {
	key   string
	value string
}

// ParseDotEnv parses .env file content into ordered KEY=VALUE pairs.
// Blank lines and # comments are skipped; an optional "export " prefix is
// accepted; values may be wrapped in single or double quotes
// (KWL-Q7X4M KWL-DOTV-002).
func ParseDotEnv(data []byte) ([]dotenvPair, error) {
	var pairs []dotenvPair
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("line %d: expected KEY=VALUE", lineNo)
		}
		key = strings.TrimSpace(key)
		value = unquoteDotEnvValue(strings.TrimSpace(value))
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNo)
		}
		pairs = append(pairs, dotenvPair{key: key, value: value})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return pairs, nil
}

// unquoteDotEnvValue strips one matching layer of single or double quotes
// and drops a trailing inline comment outside quotes for unquoted values.
func unquoteDotEnvValue(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	if i := strings.Index(v, " #"); i >= 0 {
		v = strings.TrimSpace(v[:i])
	}
	return v
}

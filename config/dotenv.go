package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DotEnvPair is one KEY=VALUE entry parsed from a .env file.
type DotEnvPair struct {
	Key   string
	Value string
}

// LoadDotEnv parses the .env file at path and exports its KEY=VALUE pairs
// into the process environment (KWL-2X1QZ CFG-DOTV-001). Variables already
// present in the process environment are never overwritten, so the real
// environment wins over .env. A missing file is not an error.
func LoadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("config: read %q: %w", path, err)
	}
	pairs, err := ParseDotEnv(data)
	if err != nil {
		return fmt.Errorf("config: parse %q: %w", path, err)
	}
	for _, kv := range pairs {
		if _, exists := os.LookupEnv(kv.Key); !exists {
			if err := os.Setenv(kv.Key, kv.Value); err != nil {
				return fmt.Errorf("config: setenv %q: %w", kv.Key, err)
			}
		}
	}
	return nil
}

// ParseDotEnv parses .env file content into ordered KEY=VALUE pairs.
// Blank lines and # comments are skipped; an optional "export " prefix is
// accepted; values may be wrapped in single or double quotes and unquoted
// values drop trailing inline comments (KWL-2X1QZ CFG-DOTV-002).
func ParseDotEnv(data []byte) ([]DotEnvPair, error) {
	var pairs []DotEnvPair
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
		pairs = append(pairs, DotEnvPair{Key: key, Value: value})
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

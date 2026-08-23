// Tests for KWN-Q7X4M
package cfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Spec: KWN-Q7X4M KWL-DOTV-002 Scope: Package
func TestKWL_DOTV_002_ParseDotEnv_ToleratesCommentsQuotesAndExport(t *testing.T) {
	body := `# comment line

PLAIN=value
export EXPORTED=yes
QUOTED_DOUBLE="hello world"
QUOTED_SINGLE='single # not comment'
UNQUOTED_COMMENTED=abc # trailing comment
`
	pairs, err := ParseDotEnv([]byte(body))
	if err != nil {
		t.Fatalf("ParseDotEnv: %v", err)
	}
	got := map[string]string{}
	for _, kv := range pairs {
		got[kv.key] = kv.value
	}
	want := map[string]string{
		"PLAIN":              "value",
		"EXPORTED":           "yes",
		"QUOTED_DOUBLE":      "hello world",
		"QUOTED_SINGLE":      "single # not comment",
		"UNQUOTED_COMMENTED": "abc",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %q = %q, want %q", k, got[k], v)
		}
	}
}

// Spec: KWN-Q7X4M KWL-DOTV-001 KWL-DOTV-002 Scope: Package
func TestKWL_DOTV_001_LoadDotEnv_SetsMissingKeepsExistingAndToleratesAbsentFile(t *testing.T) {
	t.Setenv("EXISTING", "from-shell")
	t.Setenv("SET_BY_FILE", "")
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	body := "EXISTING=from-dotenv\nNEW=loaded\nSET_BY_FILE=file\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := LoadDotEnv(path); err != nil {
		t.Fatalf("LoadDotEnv: %v", err)
	}
	if got := os.Getenv("EXISTING"); got != "from-shell" {
		t.Errorf("EXISTING = %q, want shell value to win over .env", got)
	}
	if got := os.Getenv("NEW"); got != "loaded" {
		t.Errorf("NEW = %q, want %q", got, "loaded")
	}
	// SET_BY_FILE exists but is empty in the process env; LookupEnv reports
	// existence, so .env must not overwrite it.
	if got := os.Getenv("SET_BY_FILE"); got != "" {
		t.Errorf("SET_BY_FILE = %q, want empty (existing env never overwritten)", got)
	}

	if err := LoadDotEnv(filepath.Join(dir, "missing.env")); err != nil {
		t.Errorf("LoadDotEnv(missing) = %v, want nil for absent file", err)
	}
}

// Spec: KWN-Q7X4M KWL-DOTV-003 Scope: Package
func TestKWL_DOTV_003_ParseDotEnv_MalformedLineErrorsWithLineNumber(t *testing.T) {
	_, err := ParseDotEnv([]byte("GOOD=1\nnot-a-pair\n"))
	if err == nil {
		t.Fatal("ParseDotEnv(malformed) = nil error, want error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error %q should name the offending line number 2", err.Error())
	}
}

// Spec: KWN-Q7X4M KWL-CFGV-004 Scope: Package
func TestKWL_CFGV_004_GetOr_ReturnsFallbackWhenAbsentOrEmpty(t *testing.T) {
	c := Config{"A": "x", "B": ""}
	if got := c.GetOr("A", "d"); got != "x" {
		t.Errorf("GetOr(A) = %q, want x", got)
	}
	if got := c.GetOr("B", "d"); got != "d" {
		t.Errorf("GetOr(B empty) = %q, want fallback d", got)
	}
	if got := c.GetOr("C", "d"); got != "d" {
		t.Errorf("GetOr(C absent) = %q, want fallback d", got)
	}
}

// Tests for KWL-P8W2N
package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/krewire/libs/core"
)

// Spec: KWL-P8W2N KWL-LOGV-004 Scope: Package
func TestKWL_LOGV_004_Setup_LevelAndFormat(t *testing.T) {
	if _, ok := Setup(core.EnvLocal, false).Handler().(*slog.TextHandler); !ok {
		t.Error("local handler should be text")
	}
	if _, ok := Setup(core.EnvTesting, false).Handler().(*slog.TextHandler); !ok {
		t.Error("testing handler should be text")
	}
	if _, ok := Setup(core.EnvProduction, false).Handler().(*slog.JSONHandler); !ok {
		t.Error("production handler should be JSON")
	}
	if Setup(core.EnvLocal, true).Enabled(nil, slog.LevelDebug) != true {
		t.Error("debug logger should enable Debug level")
	}
	if Setup(core.EnvProduction, false).Enabled(nil, slog.LevelDebug) {
		t.Error("production without debug should not enable Debug level")
	}
}

// Spec: KWL-P8W2N KWL-LOGV-004 Scope: Package
func TestKWL_LOGV_004_JSONHandler_StructuredRecord(t *testing.T) {
	var buf bytes.Buffer
	rec := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	rec.Info("hello", "k", "v")

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("record is not JSON: %v", err)
	}
	if m["msg"] != "hello" || m["k"] != "v" {
		t.Errorf("record fields = %v", m)
	}
}

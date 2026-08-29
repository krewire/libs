// Tests for KWL-P8W2N
package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/krewire/libs/vein"
)

// Spec: KWL-P8W2N KWL-LOGV-004 Scope: Unit
func TestKWL_LOGV_004_Setup_LevelAndFormat(t *testing.T) {
	if _, ok := Setup(vein.EnvLocal, false).Handler().(*slog.TextHandler); !ok {
		t.Error("local handler should be text")
	}
	if _, ok := Setup(vein.EnvTesting, false).Handler().(*slog.TextHandler); !ok {
		t.Error("testing handler should be text")
	}
	if _, ok := Setup(vein.EnvProduction, false).Handler().(*slog.JSONHandler); !ok {
		t.Error("production handler should be JSON")
	}
	if Setup(vein.EnvLocal, true).Enabled(nil, slog.LevelDebug) != true {
		t.Error("debug logger should enable Debug level")
	}
	if Setup(vein.EnvProduction, false).Enabled(nil, slog.LevelDebug) {
		t.Error("production without debug should not enable Debug level")
	}
}

// Spec: KWL-P8W2N KWL-LOGV-004 Scope: Unit
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

// Spec: KWL-P8W2N KWL-ERRV-008 Scope: Unit
func TestKWL_ERRV_008_LogError_CarriesAttrsAndHintIntoRecord(t *testing.T) {
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	err := vein.WithHint(
		vein.WithAttrs(vein.FailureError("cannot read config"), vein.Attr{Key: "file", Value: "krewire.yaml"}),
		"run 'kiw init' to create one",
	)
	LogError(l, "build failed", err)

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("record is not JSON: %v", err)
	}
	if m["error"] != "cannot read config" || m["file"] != "krewire.yaml" || m["hint"] != "run 'kiw init' to create one" {
		t.Errorf("record fields = %v", m)
	}
}

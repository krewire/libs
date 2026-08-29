package vein

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestSetup(t *testing.T) {
	if _, ok := Setup(EnvLocal, false).Handler().(*slog.TextHandler); !ok {
		t.Error("local handler should be text")
	}
	if _, ok := Setup(EnvTesting, false).Handler().(*slog.TextHandler); !ok {
		t.Error("testing handler should be text")
	}
	if _, ok := Setup(EnvProduction, false).Handler().(*slog.JSONHandler); !ok {
		t.Error("production handler should be JSON")
	}
	if Setup(EnvLocal, true).Enabled(nil, slog.LevelDebug) != true {
		t.Error("debug logger should enable Debug level")
	}
	if Setup(EnvProduction, false).Enabled(nil, slog.LevelDebug) {
		t.Error("production without debug should not enable Debug level")
	}
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	err := WithHint(
		WithAttrs(FailureError("cannot read config"), Attr{Key: "file", Value: "krewire.yaml"}),
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

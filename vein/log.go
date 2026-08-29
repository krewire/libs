// Package vein provides the Krewire observability foundation — logging,
// diagnostics, error handling, and stack traces (Krewire Vein).
package vein

import (
	"context"
	"log/slog"
	"os"
)

// Setup returns the process logger for env/debug, writing to stderr:
//
//   - debug on:        Debug level with source references
//   - env production:  JSON records at Info level
//   - otherwise:       text records at Info level
func Setup(env Env, debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	}
	if env == EnvProduction {
		return slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opts))
}

// Install makes Setup's logger the process default and returns it.
func Install(env Env, debug bool) *slog.Logger {
	l := Setup(env, debug)
	slog.SetDefault(l)
	return l
}

// ErrAttrs converts the diagnostic attributes attached to err into slog
// attributes so structured context travels from the error into log records
// (KWL-P8W2N KWL-ERRV-008).
func ErrAttrs(err error) []slog.Attr {
	as := AttrsOf(err)
	out := make([]slog.Attr, 0, len(as)+1)
	for _, a := range as {
		out = append(out, slog.Any(a.Key, a.Value))
	}
	return out
}

// LogError logs msg at Error level with the error's message plus any attached
// attrs and nearest hint as structured fields.
func LogError(l *slog.Logger, msg string, err error) {
	if l == nil {
		l = slog.Default()
	}
	attrs := append(ErrAttrs(err), slog.String("error", err.Error()))
	if hint := HintOf(err); hint != "" {
		attrs = append(attrs, slog.String("hint", hint))
	}
	l.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

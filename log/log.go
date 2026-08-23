// Package log provides the ecosystem's canonical logger factory: one place
// decides handler format and verbosity from the environment and debug switch
// so every binary logs consistently (KWL-P8W2N KWL-LOGV-004).
package log

import (
	"log/slog"
	"os"

	"github.com/krewire/libs/core"
)

// Setup returns the process logger for env/debug, writing to stderr:
//
//   - debug on:        Debug level with source references
//   - env production:  JSON records at Info level
//   - otherwise:       text records at Info level
func Setup(env core.Env, debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	}
	if env == core.EnvProduction {
		return slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opts))
}

// Install makes Setup's logger the process default and returns it.
func Install(env core.Env, debug bool) *slog.Logger {
	l := Setup(env, debug)
	slog.SetDefault(l)
	return l
}

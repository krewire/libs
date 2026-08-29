// Package log provides the ecosystem's canonical logger factory.
//
// Deprecated: use github.com/krewire/libs/vein (Krewire Vein) for
// observability. This package re-exports vein for backward compatibility.
package log

import "github.com/krewire/libs/vein"

var (
	Setup    = vein.Setup
	Install  = vein.Install
	ErrAttrs = vein.ErrAttrs
	LogError = vein.LogError
)

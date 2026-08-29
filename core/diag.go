// Package core provides shared primitives for the Krewire ecosystem.
//
// Deprecated: diagnostics (Attr, WithAttrs, Hint, FormatTree) have moved to
// github.com/krewire/libs/vein. This file re-exports vein for compatibility.
package core

import "github.com/krewire/libs/vein"

// Attr is one structured key/value pair attached to an error.
//
// Deprecated: use vein.Attr.
type Attr = vein.Attr

var (
	WithAttrs  = vein.WithAttrs
	AttrsOf    = vein.AttrsOf
	WithHint   = vein.WithHint
	HintOf     = vein.HintOf
	FormatTree = vein.FormatTree
)

// Package core provides shared primitives for the Krewire ecosystem.
//
// Deprecated: stack traces (WithStack, StackOf, FormatStack) have moved to
// github.com/krewire/libs/vein. This file re-exports vein for compatibility.
package core

import "github.com/krewire/libs/vein"

// StackFrame is one rendered-ready entry of a captured stack.
//
// Deprecated: use vein.StackFrame.
type StackFrame = vein.StackFrame

var (
	WithStack   = vein.WithStack
	StackOf     = vein.StackOf
	FormatStack = vein.FormatStack
)

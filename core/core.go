// Package core provides shared primitives for the Krewire ecosystem,
// including the common error type and process exit-code mapping.
//
// Deprecated: observability primitives (ExitCode, Error, WithStack, FormatTree,
// logging) have moved to github.com/krewire/libs/vein (Krewire Vein).
// This file re-exports vein for backward compatibility; new code should
// import vein directly.
package core

import "github.com/krewire/libs/vein"

// ExitCode is a standard process exit code.
//
// Deprecated: use vein.ExitCode.
type ExitCode = vein.ExitCode

const (
	ExitCodeSuccess = vein.ExitCodeSuccess
	ExitCodeFailure = vein.ExitCodeFailure
	ExitCodeUsage   = vein.ExitCodeUsage
)

// Error pairs a human-readable message with an ExitCode.
//
// Deprecated: use vein.Error.
type Error = vein.Error

var (
	NewError        = vein.NewError
	UsageError      = vein.UsageError
	FailureError    = vein.FailureError
	ExitCodeFromInt = vein.ExitCodeFromInt
)

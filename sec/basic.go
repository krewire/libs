// Package sec — BasicAuth re-export.
//
// Deprecated: use github.com/krewire/libs/auth.BasicAuth.
package sec

import "github.com/krewire/libs/auth"

// BasicVerifier validates an identifier/password pair.
//
// Deprecated: use auth.BasicVerifier.
type BasicVerifier = auth.BasicVerifier

var (
	BasicAuth = auth.BasicAuth
)

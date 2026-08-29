// Package core provides shared primitives for the Krewire ecosystem.
//
// Deprecated: environment handling will be consolidated in vein; new code
// should consider vein.Env. This file re-exports vein for compatibility.
package core

import "github.com/krewire/libs/vein"

// Env is the target environment a workload runs in.
//
// Deprecated: use vein.Env.
type Env = vein.Env

const (
	EnvLocal      = vein.EnvLocal
	EnvProduction = vein.EnvProduction
	EnvTesting    = vein.EnvTesting
	DefaultEnv    = vein.DefaultEnv
)

var (
	Envs     = vein.Envs
	ParseEnv = vein.ParseEnv
)

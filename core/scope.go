package core

import (
	"fmt"
	"strings"
)

// Scope is the ecosystem level a spec, test, or doc targets.
// Ordered: Workspace < Module < Domain < Package < Service < Func.
// See KWL-ARCH-J2K9Q.
//
// Workspace is the Krewire Workspace (hub dir ~/Workspace/Dev/krewire with bin/kiw + 7 repos),
// which is also the Go `go.work` workspace at the hub root (see `go.work` and `AGENTS.md`).
// Domain is a DDD bounded context (e.g. catalog, user) — a cohesive set of Packages
// inside a Module (internal/<domain>/). Pre-extraction it is Package set, post-extraction
// it becomes its own Service (and often its own Module/Project).
// Service is a Krewire runtime deployable, implemented in Go as a main package
// (e.g. cmd/<service>/ or service/<name>/); Func is inside Package.
type Scope string

const (
	ScopeWorkspace Scope = "Workspace"
	ScopeModule    Scope = "Module"
	ScopeDomain    Scope = "Domain"
	ScopePackage   Scope = "Package"
	ScopeService   Scope = "Service"
	ScopeFunc      Scope = "Func"
)

// AllScopes lists every valid Scope in canonical order.
var AllScopes = []Scope{ScopeWorkspace, ScopeModule, ScopeDomain, ScopePackage, ScopeService, ScopeFunc}

// scopeLevel maps Scope to its ordering index.
var scopeLevel = map[Scope]int{
	ScopeWorkspace: 0,
	ScopeModule:    1,
	ScopeDomain:    2,
	ScopePackage:   3,
	ScopeService:   4,
	ScopeFunc:      5,
}

// IsValid reports whether s is one of the six known scopes.
func (s Scope) IsValid() bool {
	_, ok := scopeLevel[s]
	return ok
}

// Level returns the ordering index (0..5). Invalid scopes return -1.
func (s Scope) Level() int {
	if v, ok := scopeLevel[s]; ok {
		return v
	}
	return -1
}

// Less reports whether s is ordered before other.
func (s Scope) Less(other Scope) bool {
	return s.Level() < other.Level()
}

// ParseScope parses s as a Scope, case-insensitive, returning UsageError on unknown.
// Accepted forms are the canonical names (Workspace, Module, Domain, Package, Service, Func),
// case-insensitive, with surrounding whitespace trimmed.
// "Project" is no longer a valid scope — use Module (Go module, formerly Project==Module).
func ParseScope(s string) (Scope, error) {
	trim := strings.TrimSpace(s)
	lower := strings.ToLower(trim)
	switch lower {
	case "workspace":
		return ScopeWorkspace, nil
	case "module":
		return ScopeModule, nil
	case "domain":
		return ScopeDomain, nil
	case "package":
		return ScopePackage, nil
	case "service":
		return ScopeService, nil
	case "func":
		return ScopeFunc, nil
	default:
		return "", UsageError(fmt.Sprintf("unknown scope %q: want one of Workspace, Module, Domain, Package, Service, Func", s))
	}
}

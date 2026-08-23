package kern

import (
	"fmt"
	"sort"

	"github.com/krewire/libs/core"
)

// Module is a kernel module. Modules are registered by name and initialized
// against the kernel. Optionally, a module may implement DependsOn() []string
// to declare ordering.
type Module interface {
	Name() string
	Init(*Kernel) error
}

// Registry holds modules by name.
type Registry struct {
	mods map[string]Module
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{mods: make(map[string]Module)}
}

// Register adds a module. Duplicate names return UsageError.
func (r *Registry) Register(m Module) error {
	if m == nil {
		return core.UsageError("module is nil")
	}
	name := m.Name()
	if name == "" {
		return core.UsageError("module name is required")
	}
	if _, exists := r.mods[name]; exists {
		return core.UsageError(fmt.Sprintf("duplicate module %q", name))
	}
	r.mods[name] = m
	return nil
}

// Resolve returns the module with the given name.
func (r *Registry) Resolve(name string) (Module, bool) {
	m, ok := r.mods[name]
	return m, ok
}

// Ordered returns modules topologically sorted by DependsOn() if implemented,
// otherwise in registration order (sorted by name for determinism).
func (r *Registry) Ordered() []Module {
	// Simple deterministic ordering: sort by name, but respect DependsOn if present.
	// For v1, we do not implement full topological sort; we sort by name and
	// ensure dependencies appear before dependents when possible via a stable sort.
	names := make([]string, 0, len(r.mods))
	for n := range r.mods {
		names = append(names, n)
	}
	sort.Strings(names)

	// If any module declares DependsOn, we attempt to reorder to satisfy it.
	// This is a best-effort bubble: repeatedly move dependents after dependencies.
	// For complex graphs, callers should register in dependency order.
	deps := make(map[string][]string)
	for _, n := range names {
		if d, ok := r.mods[n].(interface{ DependsOn() []string }); ok {
			deps[n] = d.DependsOn()
		}
	}
	// Simple insertion sort respecting deps: if a appears before its dependency, move it after.
	for i := 0; i < len(names); i++ {
		for _, dep := range deps[names[i]] {
			// Find dep index
			depIdx := -1
			for j, n := range names {
				if n == dep {
					depIdx = j
					break
				}
			}
			if depIdx >= 0 && depIdx > i {
				// Move dependency before dependent
				depName := names[depIdx]
				// Remove from depIdx
				names = append(names[:depIdx], names[depIdx+1:]...)
				// Insert at i
				names = append(names[:i], append([]string{depName}, names[i:]...)...)
				// Restart scan
				i = -1
				break
			}
		}
	}
	out := make([]Module, 0, len(names))
	for _, n := range names {
		out = append(out, r.mods[n])
	}
	return out
}

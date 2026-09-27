package kern

import (
	"fmt"
	"sort"
	"sync"

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
	mu   sync.RWMutex
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
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.mods[name]; exists {
		return core.UsageError(fmt.Sprintf("duplicate module %q", name))
	}
	r.mods[name] = m
	return nil
}

// Resolve returns the module with the given name.
func (r *Registry) Resolve(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mods[name]
	return m, ok
}

// Ordered returns modules topologically sorted by DependsOn() if implemented,
// otherwise in registration order (sorted by name for determinism).
// Circular dependencies are resolved safely without infinite loops.
func (r *Registry) Ordered() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.mods))
	for n := range r.mods {
		names = append(names, n)
	}
	sort.Strings(names)

	deps := make(map[string][]string, len(names))
	for _, n := range names {
		if d, ok := r.mods[n].(interface{ DependsOn() []string }); ok {
			deps[n] = d.DependsOn()
		}
	}

	inDegree := make(map[string]int, len(names))
	dependents := make(map[string][]string, len(names))
	for _, n := range names {
		inDegree[n] = 0
	}
	for n, ds := range deps {
		for _, dep := range ds {
			if _, exists := r.mods[dep]; exists {
				inDegree[n]++
				dependents[dep] = append(dependents[dep], n)
			}
		}
	}

	var queue []string
	for _, n := range names {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}

	var sorted []string
	visited := make(map[string]bool, len(names))

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if visited[curr] {
			continue
		}
		visited[curr] = true
		sorted = append(sorted, curr)

		var nextCandidates []string
		for _, dep := range dependents[curr] {
			inDegree[dep]--
			if inDegree[dep] == 0 && !visited[dep] {
				nextCandidates = append(nextCandidates, dep)
			}
		}
		sort.Strings(nextCandidates)
		queue = append(queue, nextCandidates...)
	}

	// If there were circular dependencies, append any remaining nodes deterministically
	// to prevent infinite loops.
	if len(sorted) < len(names) {
		for _, n := range names {
			if !visited[n] {
				sorted = append(sorted, n)
			}
		}
	}

	out := make([]Module, 0, len(sorted))
	for _, n := range sorted {
		out = append(out, r.mods[n])
	}
	return out
}

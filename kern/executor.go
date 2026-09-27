package kern

import (
	"context"
	"sync"

	"github.com/krewire/libs/core"
	"github.com/krewire/libs/vein"
)

// Executor dispatches a workload.
type Executor interface {
	Execute(ctx context.Context, workload core.Workload) vein.ExitCode
}

// executor dispatches to the module that handles the workload's Kind.
type executor struct {
	mu       sync.RWMutex
	handlers map[core.Kind]func(context.Context, core.Workload) vein.ExitCode
}

func newExecutor() *executor {
	return &executor{handlers: make(map[core.Kind]func(context.Context, core.Workload) vein.ExitCode)}
}

func (e *executor) Register(kind core.Kind, fn func(context.Context, core.Workload) vein.ExitCode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[kind] = fn
}

func (e *executor) Execute(ctx context.Context, workload core.Workload) vein.ExitCode {
	e.mu.RLock()
	fn, ok := e.handlers[workload.Kind]
	e.mu.RUnlock()
	if ok {
		return fn(ctx, workload)
	}
	return vein.ExitCodeUsage
}

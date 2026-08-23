package kern

import (
	"context"

	"github.com/krewire/libs/core"
)

// Executor dispatches a workload.
type Executor interface {
	Execute(ctx context.Context, workload core.Workload) core.ExitCode
}

// executor dispatches to the module that handles the workload's Kind.
type executor struct {
	handlers map[core.Kind]func(context.Context, core.Workload) core.ExitCode
}

func newExecutor() *executor {
	return &executor{handlers: make(map[core.Kind]func(context.Context, core.Workload) core.ExitCode)}
}

func (e *executor) Register(kind core.Kind, fn func(context.Context, core.Workload) core.ExitCode) {
	e.handlers[kind] = fn
}

func (e *executor) Execute(ctx context.Context, workload core.Workload) core.ExitCode {
	if fn, ok := e.handlers[workload.Kind]; ok {
		return fn(ctx, workload)
	}
	return core.ExitCodeUsage
}

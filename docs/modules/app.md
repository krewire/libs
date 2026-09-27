# `app`

Import: `github.com/krewire/libs/app`

## Purpose

`app` provides application bootstrap, lifecycle coordination, and a named service container. It coordinates providers without owning their concrete implementations.

## Main API

- `Application` — provider registration and lifecycle coordination.
- `Runner` — runtime contract supplied by the [`runner`](./runner.md) module.
- `RunnerFunc` — compatibility alias for `runner.Func`.
- `New()` — creates an application with an empty container.
- `Use(providers...)` — registers providers in startup order.
- `Bootstrap(ctx)` — registers all providers, then starts lifecycle-aware providers.
- `Shutdown(ctx)` — stops started providers in reverse order.
- `Run(ctx, runner)` — bootstraps providers, runs the supplied runtime, and shuts down providers.
- `Container` — named service registry.
- `Resolve[T]` — typed service lookup.

## Example

```go
application := app.New()
if err := application.Use(databaseProvider, httpProvider); err != nil { return err }
if err := application.Run(ctx, app.RunnerFunc(func(ctx context.Context, services *app.Container) error {
    return runCLIOrServer(ctx, services)
})); err != nil { return err }
```

## Design boundary

`app` owns orchestration and dependency lookup. Concrete service behavior belongs in [`service`](./service.md) providers. Workload execution remains in [`kern`](./kern.md); `app` does not replace the kernel.

# `app`

Import: `github.com/krewire/libs/app`

## Purpose

`app` provides application bootstrap, lifecycle coordination, and a named service container. It coordinates providers without owning their concrete implementations.

## Main API

- `Application` — provider registration and bootstrap/shutdown lifecycle.
- `New()` — creates an application with an empty container.
- `Use(providers...)` — registers providers in startup order.
- `Bootstrap(ctx)` — registers all providers, then starts lifecycle-aware providers.
- `Shutdown(ctx)` — stops started providers in reverse order.
- `Container` — named service registry.
- `Resolve[T]` — typed service lookup.

## Example

```go
application := app.New()
if err := application.Use(databaseProvider, httpProvider); err != nil { return err }
if err := application.Bootstrap(ctx); err != nil { return err }
defer application.Shutdown(context.Background())
```

## Design boundary

`app` owns orchestration and dependency lookup. Concrete service behavior belongs in [`service`](./service.md) providers. Workload execution remains in [`kern`](./kern.md); `app` does not replace the kernel.

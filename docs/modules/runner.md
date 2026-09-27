# `runner`

Import: `github.com/krewire/libs/runner`

## Purpose

`runner` defines the runtime contract for all Krewire workloads. It does not assume an HTTP server, CLI process, worker, scheduler, event consumer, or other application type.

## Main API

- `Runner` — executes a runtime with a `context.Context` and `service.Registry`.
- `Func` — adapts a function into a `Runner`.

## Example

```go
application := app.New()
_ = application.Use(databaseProvider)

return application.Run(ctx, runner.Func(func(
    ctx context.Context,
    services service.Registry,
) error {
    return runWorkload(ctx, services)
}))
```

## Design boundary

`runner` defines execution shape only. Application bootstrap and lifecycle belong to [`app`](./app.md); provider contracts belong to [`service`](./service.md); kernel workload orchestration belongs to [`kern`](./kern.md).

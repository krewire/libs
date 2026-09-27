# Krewire Libraries

Reusable Go libraries for the Krewire ecosystem: small, composable packages for domain rules, application lifecycle, configuration, validation, security, observability, terminal output, and Markdown rendering.

Module: `github.com/krewire/libs`

## What this repository provides

The repository is a Go monorepo. Each top-level package can be imported independently; consumers should depend only on the packages they use.

| Package | Purpose | Main API areas |
|---|---|---|
| [`core`](./core) | Shared domain model and invariants | `Kind`, `Workload`, `Project`, `Scope`, `SpecID`, `RequirementID`, `DomainEvent`, versions, exit codes, diagnostic/error aliases |
| [`kern`](./kern) | Kernel boot and workload execution | `Kernel`, `Module`, `Registry`, `Executor`, `Supervisor` |
| [`vein`](./vein) | Logging, diagnostics, errors, and stack traces | `Setup`, `Install`, `WithAttrs`, `WithHint`, `FormatTree`, `WithStack`, `StackOf` |
| [`auth`](./auth) | HTTP authentication primitives | `BasicAuth`, `JWTAuth`, `SignJWT`, `ParseJWT`, `Identity`, HTTP errors |
| [`sec`](./sec) | HTTP security middleware and PII-safe diagnostics | security headers, CORS, CSRF, health endpoints, authentication aliases, policies, SSRF URL validation, `MaskPII`, `MaskPIIMap`, `MaskPIIAttrs` |
| [`config`](./config) | Typed configuration and environment overlays | `Load`, `LoadOrDefault`, `Override`, `.env` parsing, `Vars` |
| [`validation`](./validation) | Reflection-based struct validation with extensible rule evaluators | `Struct`, `Field`, `rules.Register`, tags such as `required`, `email`, `min`, `max`, `len`, `oneof`, `pattern` |
| [`term`](./term) | Terminal detection and ANSI styling | `Terminal`, `NewTerminal`, `Paint`, colors and styles |
| [`markdown`](./markdown) | Markdown-to-HTML rendering | GFM rendering, heading IDs, base-prefixed links |

`core` is the declarative center of the ecosystem: it defines valid project kinds, workloads, scopes, versions, and project invariants. `kern` is the imperative layer: it initializes modules, orders dependencies, boots a kernel, dispatches workloads, and coordinates shutdown. `vein` is the preferred observability package; older integrations may still expose compatibility aliases through `core` and `sec`.

## Supported domain concepts

`core.Kind` recognizes these project kinds:

`app`, `cli`, `site`, `book`, `worker`, `service`, `infra`, and `kernel`.

`core.Scope` models the hierarchy `workspace → module → domain → service → unit`. `core.Version` and the compatibility helpers support checking module and ecosystem version requirements before execution.

## Installation

Use the module or package you need:

```bash
go get github.com/krewire/libs/core
go get github.com/krewire/libs/config
go get github.com/krewire/libs/sec
```

The module currently targets Go `1.26.0` and depends on:

- `github.com/yuin/goldmark` for Markdown rendering
- `gopkg.in/yaml.v3` for YAML configuration and variable storage
- `golang.org/x/net` as an indirect dependency

## Small examples

### Load typed configuration with environment overrides

```go
type ServerConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

var cfg ServerConfig
if err := config.Load("krewire.yaml", &cfg); err != nil {
    return err
}
if err := config.Override(&cfg, os.LookupEnv, config.WithPrefix("APP_")); err != nil {
    return err
}
```

### Validate a struct

```go
type User struct {
    Email string `validate:"required,email"`
    Role  string `validate:"oneof=admin viewer"`
}

if err := validation.Struct(User{Email: "user@example.com", Role: "viewer"}); err != nil {
    // *validation.ValidationError contains field-level failures.
    return err
}
```

### Protect an HTTP handler

```go
handler := sec.SecurityHeaders()(
    sec.CORS(sec.WithOrigins("https://example.com"))(
        sec.Require(sec.Authenticated())(http.HandlerFunc(handleRequest)),
    ),
)
```

For JWT or Basic authentication, use the middleware in [`auth`](./auth) or the corresponding aliases in [`sec`](./sec). Authentication secrets must come from runtime configuration; do not hard-code them in source code.

### Render Markdown

```go
html, err := markdown.RenderWithBase(source, "/docs/")
if err != nil {
    return err
}
```

## Development

Prerequisites:

- Go `1.26.0` or a compatible newer toolchain

Run the same checks used by CI:

```bash
gofmt -l .                 # must print nothing
go vet ./...
go test ./...
```

A complete local build can be checked with:

```bash
go build ./...
```

CI runs on pushes to `main` and pull requests. It executes formatting, `go vet`, and the complete test suite with the stable Go toolchain.

## Design principles

- Keep packages narrow and composable.
- Prefer the standard library where it already solves the problem (`net/http`, `log/slog`, `os`, `flag`).
- Keep domain rules in `core` and lifecycle mechanics in `kern`.
- Make security behavior explicit and fail safely.
- Preserve specification-to-test traceability; formal specifications live in [`docs/specs`](./docs/specs).
- Update tests and documentation when public behavior changes.

More context is available in [`docs/architecture.md`](./docs/architecture.md), [`docs/modules`](./docs/modules/README.md), [`docs/philosophy.md`](./docs/philosophy.md), and [`docs/index.md`](./docs/index.md).

## Related repositories

- [`framework`](https://github.com/krewire/framework) — unified Krewire framework
- [`mdbind`](https://github.com/krewire/mdbind) — book/site builder
- [`kiw`](https://github.com/krewire/kiw) — Krewire development tool

## License

MIT — see [`LICENSE`](./LICENSE).

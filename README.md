# Krewire Libraries

**Krewire Libraries** is a monorepo of modular, reusable Go libraries for the Krewire ecosystem — `github.com/krewire/libs`. It provides the shared building blocks behind [`framework`](https://github.com/krewire/framework), [`mdbind`](https://github.com/krewire/mdbind), and [`kiw`](https://github.com/krewire/kiw).

The unified framework vision ([`KWF-M8K2Q`](../framework/docs/specs/KWF-ARCH-M8K2Q-unified-framework-vision.md)) relies on `libs` for a single, typed `krewire.yaml` across all eight project kinds (`app`, `cli`, `site`, `book`, `worker`, `service`, `infra`, `kernel`).

## Overview

Each concern is an independent package, versioned and released on its own schedule. Consumers pull only what they need. Where the Go standard library already covers a concern (`flag`, `log/slog`, `os`), `libs` does not re-implement it.

## Packages

| Path | Description |
|------|-------------|
| `core/` | **Business rules & workload registry** — `Kind`/`Workload` matrix, `SpecID`/`RequirementID`, `Project` invariants, `DomainEvent` + shared primitives (`ExitCodeSuccess/Failure/Usage`). Declarative control plane (`KWL-K1N2Q`). |
| `kern/` | **Kernel executor & supervisor** — generic `Kernel`/`Module`/`Registry`/`Executor`/`Supervisor` for boot, lifecycle, and workload dispatch. Imperative control plane (`KWL-KERN-X8P3L`). |
| `log/` | Canonical `slog` logger factory — installs the process logger (JSON/text, level) from `core.Env`/debug; bridges error diagnostics into structured attrs. |
| `term/` | Terminal I/O, output formatting, and color conventions. |
| `config/` | Typed `krewire.yaml` loading for all 8 kinds (delegates business validation to `core`). |
| `validate/` | Struct validation (`validate:"required"` tags) for config and resource schemas. |

Together `core` (declarative) + `kern` (imperative) form the **central control** of the ecosystem: every repo imports `core` for types/rules, and `framework`/`krewire` compose via `kern`. `config` + `validate` enforce the single-config rule — every workload from `cli` to `infra` is described in one `krewire.yaml` and validated before `kiw build` / `kiw deploy`.

Standard library responsibilities that `libs` intentionally does not duplicate:

| Concern | Stdlib package |
|---------|----------------|
| Argument parsing | `flag` |
| Structured logging | `log/slog` |
| Environment / config sources | `os`, `strconv` |

## Getting Started

### Prerequisites

- Go 1.22+ — https://go.dev/dl/

### Building and testing

```bash
go build ./...
go test ./...
gofmt -l . && go vet ./...
```

## Specifications

- `KWL-CORE-K1N2Q` — Core Business Rules & Workload Registry (declarative control plane)
- `KWL-KERN-X8P3L` — Kernel Executor & Supervisor (imperative control plane)
- `KWL-CONFIG-2X1QZ` — Configuration loading
- `KWL-VALIDATE-LHANF` — Struct validation
- `KWL-CORE-W0J2X` — Errors & exit codes (extended by K1N2Q)
- `KWL-TERM-R934Y` — Terminal I/O & rendering

All specs live in `docs/specs/` (`KWL-*`).

## Related Repositories

- [framework](https://github.com/krewire/framework) — unified framework (`tui`/`web`+`ssg`/`ui`/`app`/`runtime`/`worker`/`service`/`infra`)
- [mdbind](https://github.com/krewire/mdbind) — book/site builder on `framework/web`
- [kiw](https://github.com/krewire/kiw) — devtool CLI for all kinds

## License

MIT — see [LICENSE](LICENSE).

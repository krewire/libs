# Krewire Vein — Observability

**Krewire Vein** (`github.com/krewire/libs/vein`) is the observability subproduct of `libs` — logging, diagnostics, error handling, and stack traces. Like Spring Boot Actuator for observability, Vein is the single place that decides how every Krewire binary logs, reports errors, and traces failures.

## Packages

| Symbol | Description |
|--------|-------------|
| `ExitCode`, `Error`, `UsageError`, `FailureError` | Process exit codes 0/1/2 and typed errors carrying them |
| `Env`, `ParseEnv`, `Envs` | Environment (`local`/`production`/`testing`) for log format selection |
| `Attr`, `WithAttrs`, `AttrsOf`, `WithHint`, `HintOf`, `FormatTree` | Structured diagnostics attached to errors without breaking `errors.Is/As` |
| `WithStack`, `StackOf`, `FormatStack`, `StackFrame` | Immutable stack capture and rendering |
| `Setup`, `Install`, `ErrAttrs`, `LogError` | Canonical `log/slog` factory (JSON in production, text elsewhere; `Debug` adds source) |

## Usage

```go
import "github.com/krewire/libs/vein"

err := vein.WithStack(vein.FailureError("cannot read config"))
err = vein.WithAttrs(err, vein.Attr{Key: "file", Value: "krewire.yaml"})
err = vein.WithHint(err, "run 'kiw new' first")
vein.LogError(nil, "startup failed", err) // logs with attrs + hint
fmt.Println(vein.FormatTree(err))
```

`core` and `log` re-export Vein for backward compatibility; new code should import `vein` directly.

## Specs

* `KWL-ERR-P8W2N` — error handling, stack traces, and logging
* `KWL-CORE-W0J2X` — errors & exit codes
* `KWL-ENV-K4T7W` — environments (consumed by `Setup`)

All specs live in `../docs/specs/` (`KWL-*`).

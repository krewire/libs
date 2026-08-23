# Specification — Error Handling, Stack Traces & Logging

| Field       | Value                                            |
| ----------- | ------------------------------------------------ |
| SpecID      | KWL-P8W2N                                        |
| Title       | Error Handling, Stack Traces & Structured Logging |
| Status      | Draft                                            |
| Date        | 2026-08-23                                       |
| Author      | Krewire Contributors                              |
| Domain      | Core — Diagnostics                               |

## 1. Context

Failures today surface as single-line messages: `fail()` prints `kiw: <err>`,
panics are recovered without caller context, and every binary configures
`slog.Default()` ad-hoc regardless of environment. Debugging production
incidents requires knowing *where* an error was created and *which* code path
logged it — information the current pipeline discards.

## 2. Problem Statement

Errors carry no origin metadata, so support sessions start from guesswork;
panics are logged as bare values without stacks; log format and verbosity do
not follow the environment (JSON vs text, Info vs Debug), and the web layer
ships without recovery or request logging enabled by default.

## 3. Goals

- G1 — Errors can carry an immutable, captured stack at creation/wrap point
  without breaking `errors.Is/As`.
- G2 — One canonical logger factory driven by `core.Env` + debug
  (KWL-K4T7W): text locally, JSON in production, Debug level + source when
  debugging. Lives at `github.com/krewire/libs/log` (`package log`).
- G3 — The web layer recovers panics with full stacks server-side and never
  leaks internals to clients; request logging ships enabled.
- G4 — The devtool adopts the same logger and surfaces attached stacks on
  failure when debugging.

## 4. Non-Goals

- NG1 — Distributed tracing export (OTLP spans); hooks exist via
  `framework/app.Tracer`, transport comes later.
- NG2 — Log shipping/aggregation integrations.
- NG3 — Retrofitting stacks onto every existing error construction site;
  adoption is incremental at failure boundaries.
- NG4 — Custom log destinations/files; stderr only (12-factor).

## 5. Requirements

| ID           | Requirement                                                                                     | Priority | Scope    |
| ------------ | ----------------------------------------------------------------------------------------------- | -------- | -------- |
| KWL-ERRV-001 | `core.WithStack(err)` captures the calling goroutine's PCs at wrap time; `errors.Is/As` and `Unwrap` pass through unchanged. | Must | Package |
| KWL-ERRV-002 | `core.StackOf(err)` extracts frames via `errors.As`; nil when absent. `core.FormatStack` renders `func file:line` lines. | Must | Package |
| KWL-ERRV-003 | Stacking the same error twice keeps both traces extractable.                                    | Should   | Package  |
| KWL-LOGV-004 | `logging.Setup(env, debug)` returns a stderr logger: Debug level with source when debug; JSON records in production; text otherwise. | Must | Package |
| KWF-HTTPV-005| `RecoverMiddleware` logs the panic value together with a captured goroutine stack; clients receive a 500 without internals. | Must | Package |
| KWF-HTTPV-006| `App` attaches recovery and access-log middleware by default; explicit `Use` composes outside them. | Must     | Package  |
| KWL-DIAGV-007| `bootRuntime` installs the environment-appropriate default logger; the devtool prints extracted stacks on failure when debug is on. | Must | Domain   |

## 6. Non-Functional Requirements

- NFR1 — `libs/core` stays stdlib-only (`runtime` permitted).
- NFR2 — Zero allocation on the happy path of middleware beyond one wrapper
  writer per request.
- NFR3 — Quality gates pass in every touched repo (`gofmt`, `vet`, `test`).

## 7. Success Criteria

- S1 — A panicking route returns HTTP 500 with body `internal server error`
  while the server log shows the panic plus a stack naming the handler frame.
- S2 — `kiw run --debug` on a failing project logs at Debug with source
  references; without `--debug` nothing below Info appears.
- S3 — `errors.Is(wrapped, target)` remains true through `WithStack`.

## 8. Related Specifications

| SpecID    | Title                                    |
| --------- | ---------------------------------------- |
| KWL-M1ZKS | Krewire Libraries — Initial Specification |
| KWL-K4T7W | Environments & Debug Mode                |
| KWL-W0J2X | Core Errors & Exit Codes                 |
| KWN-6K41E | krewire run/dev/deploy                    |

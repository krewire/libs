# Krewire Security — Security

**Krewire Security** (`github.com/krewire/libs/sec`) is the framework-agnostic HTTP security package for authentication integration, authorization, browser hardening, CSRF, CORS, health probes, SSRF URL validation, and PII-safe diagnostics.

## Features

| Symbol | Description |
|--------|-------------|
| `Identity`, `IdentityFrom`, `HasRole` | Authenticated caller (`Subject`, `Method`, `Roles`, `Claims`) stored in `context.Context` |
| `BasicAuth`, `JWTAuth`, `SignJWT`, `ParseJWT`, `Claims` | RFC 7617 Basic and HS256 JWT (header `Authorization: Basic/Bearer` or cookie) |
| `SecurityHeaders`, `StripTags` | Browser hardening (`X-Content-Type-Options`, `X-Frame-Options`, `CSP`, `HSTS`) |
| `CSRF`, `CSRFFrom` | Double-submit token (`XSRF-TOKEN` cookie + `X-CSRF-Token` header/form) |
| `Policy`, `Require`, `PolicySet`, `Authenticated`, `WithRoles` | Before-gate policies (`401`/`403` via `HTTPError`) |
| `CORS`, `WithOrigins`, `WithCredentials` | Explicit origin policy with wildcard-plus-credentials protection |
| `Health`, `WithLivenessPath`, `WithReadinessPath` | Liveness/readiness endpoints with configurable checks |
| `ValidateOutboundURL` | Allowlisted HTTP(S) URL validation with literal private-address rejection |
| `MaskPII`, `MaskPIIMap`, `MaskPIIAttrs` | Mask email, phone, card, IPv4/IPv6 data and redact secret-bearing keys |
| `NewPIIHandler` | `log/slog.Handler` wrapper that masks PII in records and groups |
| `HTTPError`, `Unauthorized`, `Forbidden`, `Middleware` | Structured HTTP errors and `func(http.Handler) http.Handler` middleware |

## Usage

```go
import "github.com/krewire/libs/sec"

mux := http.NewServeMux()
mux.Handle("/", sec.SecurityHeaders()(sec.CSRF()(myHandler)))
mux.Handle("/api", sec.Require(sec.Authenticated(), sec.WithRoles("admin"))(apiHandler))

// Basic
mux.Handle("/admin", sec.BasicAuth("admin", verify)(adminHandler))
// JWT
mux.Handle("/secure", sec.JWTAuth(secret)(secureHandler))
```

`framework/web` may re-export `sec` for backward compatibility; new code can import `sec` directly.

## Security boundaries

`ValidateOutboundURL` does not perform network requests and cannot prevent DNS rebinding by itself. Consumers that make outbound requests must resolve and validate destination addresses safely, restrict ports, and connect using the validated resolution. `StripTags` extracts text; it is not a context-safe HTML sanitizer. Use proper output encoding or a dedicated sanitizer before rendering untrusted HTML.

PII masking is for logs and diagnostics. It is not authorization, encryption, or a data-retention control.

## Specs

- [`KWL-O10SC`](../docs/specs/KWL-SEC-O10SC-owasp-top-10-security-controls.md) — OWASP Top 10 security controls
- [`KWL-ERR-P8W2N`](../docs/specs/KWL-ERR-P8W2N-error-handling-stack-traces-and-logging.md) — error and logging behavior

All repository specifications live in [`../docs/specs/`](../docs/specs/).

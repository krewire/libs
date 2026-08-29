# Krewire Security — Security

**Krewire Security** (`github.com/krewire/libs/sec`) is the security subproduct of `libs` — authentication, authorization, and browser hardening. Like Spring Security, it is framework-agnostic and works with `net/http` handlers.

## Features

| Symbol | Description |
|--------|-------------|
| `Identity`, `IdentityFrom`, `HasRole` | Authenticated caller (`Subject`, `Method`, `Roles`, `Claims`) stored in `context.Context` |
| `BasicAuth`, `JWTAuth`, `SignJWT`, `ParseJWT`, `Claims` | RFC 7617 Basic and HS256 JWT (header `Authorization: Basic/Bearer` or cookie) |
| `SecurityHeaders`, `StripTags` | Browser hardening (`X-Content-Type-Options`, `X-Frame-Options`, `CSP`, `HSTS`) |
| `CSRF`, `CSRFFrom` | Double-submit token (`XSRF-TOKEN` cookie + `X-CSRF-Token` header/form) |
| `Policy`, `Require`, `PolicySet`, `Authenticated`, `WithRoles` | Before-gate policies (`401`/`403` via `HTTPError`) |
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

`framework/web` re-exports `sec` for backward compatibility; new code can import `sec` directly.

## Specs

Security behavior is specified in `framework` (`KWF-WEB-R9T4C`, `KWF-WEB-B2X7D`) and consumed via `sec`.

# Specification — Core Business Rules & Workload Registry

| Field  | Value                          |
|--------|--------------------------------|
| SpecID | KWL-K1N2Q                      |
| Title  | Core Business Rules & Workload Registry |
| Status | Draft                          |
| Date   | 2026-08-21                     |
| Author | Krewire Contributors            |
| Domain | Libraries — Core               |

## 1. Context

`libs/core` currently exposes only process primitives (`ExitCode`, `Error` — `KWL-W0J2X`). Business rules that define the ecosystem — the 8 `project.kind` values, the 9-workload matrix, `SpecID`/`RequirementID` formats, and invariants (*one `krewire.yaml` only*, *spec-first*, *opt-in batteries*, *no `go.work`*) — are scattered across `AGENTS.md`, `krewire/internal/shape`, and `KWF-M8K2Q`. Without a single source, validation is duplicated and drift is possible. The unified vision needs `libs/core` to be the declarative control plane: the domain layer every other package imports for types and validation.

## 2. Problem Statement

- Kind/workload definitions live in multiple places; adding `worker`/`service`/`infra` required touching 10+ files.
- SpecID and requirement ID formats are documented but not enforced as types.
- Invariants are enforced ad-hoc in `krewire` and `framework`, not as reusable business rules.
- `libs/config` and `libs/validate` cannot delegate workload-aware validation to a shared authority.

## 3. Goals

- G1 — `libs/core` becomes the single authority for `Kind`, `Workload`, `SpecID`, `RequirementID`, `Project`, and domain events.
- G2 — 100% backward compatible: existing `ExitCode`/`Error` API unchanged; new types are additive.
- G3 — Pure domain layer: `core` has zero dependencies outside stdlib; it never imports `framework` or `krewire`.
- G4 — All business validation is reusable via `core` so `libs/config`, `libs/validate`, `framework`, and `krewire` converge.

## 4. Non-Goals

- NG1 — Not an executor or lifecycle manager; that is `libs/kern` (`KWL-KERN-X8P3L`).
- NG2 — Not re-implementing `libs/config` (YAML loading) or `libs/validate` (struct validation); `core` provides types and business predicates those packages call.
- NG3 — No I/O, no filesystem, no `os.Getenv` inside `core`; side effects stay in `kern`/`krewire`.

## 5. Requirements

### 5.1 Kind & Workload Registry

The 8 kinds from the unified vision (`internal/docs/project-vision.md`, `KWF-M8K2Q`):

| ID          | Requirement | Priority |
|-------------|-------------|----------|
| KWL-CORE-001 | `type Kind string` with constants `KindApp`, `KindCLI`, `KindSite`, `KindBook`, `KindWorker`, `KindService`, `KindInfra`, `KindKernel`; `func (Kind) IsValid() bool`; `func ParseKind(string) (Kind, error)` returning `UsageError` on unknown. | Must |
| KWL-CORE-002 | `type Workload struct { Kind Kind; Package string; Title string; SpecID SpecID; Status Status }` and `var Workloads []Workload` covering 9 workloads (see `internal/docs/project-vision.md` table) + `func WorkloadFor(Kind) (Workload, bool)`. | Must |
| KWL-CORE-003 | `type Status string` with `StatusShipped`, `StatusPlanned`; workload statuses sourced from `KWF-M8K2Q` and kept in sync. | Must |

### 5.2 Spec & Requirement IDs

| ID          | Requirement | Priority |
|-------------|-------------|----------|
| KWL-CORE-010 | `type SpecID string` with format `{ProjectId}-{Scope}-{5-char}-{slug}`; `func ParseSpecID(string) (SpecID, error)` validates `ProjectId` in `{KWF,KWL,KWM,KWN,KWG,KWI,KWD}`, `Scope` as `[A-Z0-9]+`, 5-char alphanumeric. | Must |
| KWL-CORE-011 | `type RequirementID string` with format `FRK-*`/`KWL-*`/`KWM-*` etc. + `func ParseRequirementID(string) error`. | Should |
| KWL-CORE-012 | Helpers `SpecID.Project()`, `.Scope()`, `.Code()` for indexing. | Should |

### 5.3 Project & Invariants

| ID          | Requirement | Priority |
|-------------|-------------|----------|
| KWL-CORE-020 | `type Project struct { Name, ModulePath string; Kind Kind; ConfigPath string }` with `func (Project) Validate() error` enforcing: name kebab-case (`^[a-z][a-z0-9-]*$`), module path is import path, kind is valid. | Must |
| KWL-CORE-021 | `func ValidateKrewireYamlPath(path string) error` ensures config path is `krewire.yaml` (no `ssg.yaml`) — encodes invariant *one config file only*. | Must |
| KWL-CORE-022 | Predicate `IsOptIn(kind Kind, imported []string) bool` helper to detect opt-in cost violations (monolith importing `service`/`infra`). | Should |

### 5.4 Domain Events & Exit Integration

| ID          | Requirement | Priority |
|-------------|-------------|----------|
| KWL-CORE-030 | `type DomainEvent struct { Type string; Payload any; At time.Time }` for cross-module communication (e.g., `worker.job.enqueued`). | Should |
| KWL-CORE-031 | `Error` and `ExitCode` remain the canonical process primitives; new types use them (e.g., `ParseKind` returns `*Error` with `ExitCodeUsage`). | Must |

## 6. Non-Functional Requirements

- NFR1 — Zero breaking changes; `go vet` and `go test ./...` pass in `libs` and downstream `framework`/`krewire` unchanged.
- NFR2 — `core` stays stdlib-only; no `gopkg.in/yaml.v3` or other deps.
- NFR3 — 100% `gofmt` clean, idiomatic Go, `go doc` comments on exported types.

## 7. Success Criteria

- S1 — `go doc github.com/krewire/libs/core` lists `Kind`, `Workload`, `SpecID`, `Project` with examples; `libs/core` tests cover `ParseKind`, `ParseSpecID`, `Project.Validate` including error cases returning `ExitCodeUsage`.
- S2 — `libs/config` can delegate `Kind` validation to `core.ParseKind` (no duplicated logic).
- S3 — `krewire/internal/shape` can be refactored to use `core.Kind` without behavior change (follow-up, not in this spec).

## 8. Related Specifications

| SpecID | Title |
|--------|-------|
| [KWL-W0J2X](./KWL-CORE-W0J2X-errors-exit-codes.md) | Errors & Exit Codes (extends) |
| [KWL-2X1QZ](./KWL-CONFIG-2X1QZ-configuration-loading.md) | Configuration Loading (consumer) |
| [KWL-LHANF](./KWL-VALIDATE-LHANF-struct-validation.md) | Struct Validation (consumer) |
| [KWF-M8K2Q](../framework/KWF-ARCH-M8K2Q-unified-framework-vision.md) | Unified Vision (source of workload matrix) |
| [KWL-KERN-X8P3L](./KWL-KERN-X8P3L-kernel-executor.md) | Kernel Executor (imperative counterpart) |

## 9. References

- `internal/docs/project-vision.md` — 9-workload matrix
- `AGENTS.md` — 8-kind detection table

# `service`

Import: `github.com/krewire/libs/service`

## Purpose

`service` defines the provider contracts used by [`app`](./app.md). It contains interfaces only and does not own application state or startup orchestration.

## Main API

- `Provider` — requires `Name()` and `Register(Registry)`.
- `Registry` — exposes named `Set` and `Get` operations.
- `Starter` — optional post-registration startup hook.
- `Stopper` — optional graceful cleanup hook.

Providers should register dependencies during `Register`, perform runtime startup in `Start`, and release resources in `Stop`.

## Example

```go
type CacheProvider struct { cache *Cache }
func (p *CacheProvider) Name() string { return "cache" }
func (p *CacheProvider) Register(r service.Registry) error {
    p.cache = NewCache()
    return r.Set("cache", p.cache)
}
func (p *CacheProvider) Start(context.Context, service.Registry) error { return p.cache.Open() }
func (p *CacheProvider) Stop(context.Context, service.Registry) error { return p.cache.Close() }
```

## Design boundary

`service` is intentionally independent of `app`. This keeps provider contracts reusable and prevents providers from depending on the application implementation.

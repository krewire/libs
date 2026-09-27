# `functional`

Import: `github.com/krewire/libs/functional`

## Purpose

`functional` provides small, generic utilities built around closures. It focuses on captured configuration, captured state, function composition, and memoization without forcing callers to define helper structs.

## Main API

- `Memoize[K, V]` — caches successful function results by comparable key and is safe for concurrent use.
- `Bind[A, B, C]` — fixes the first argument of a two-argument function.
- `Compose[A, B, C]` — chains two functions.
- `Counter` — returns a stateful incrementing closure.
- `ErrNilFunction` — returned by `Memoize` when its function argument is nil.

## Example

```go
lookup := functional.Memoize(func(key string) (Record, error) {
    return repository.Find(key)
})

normalize := functional.Compose(strings.TrimSpace, strings.ToLower)
nextID := functional.Counter(100)
```

## Design boundary

This package does not replace language-level closures or standard-library primitives such as `sync.OnceValue`. It provides reusable patterns where captured state or function composition improves clarity.

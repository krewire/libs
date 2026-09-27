// Package functional provides small, generic function-composition utilities.
// The utilities return closures that capture configuration or state without
// requiring callers to define a dedicated struct.
package functional

import (
	"errors"
	"sync"
)

// ErrNilFunction indicates that a required function argument was nil.
var ErrNilFunction = errors.New("functional: nil function")

// Memoize returns a function that caches successful results by key. The
// function is safe for concurrent use. Errors are returned to the caller and
// are not cached.
func Memoize[K comparable, V any](fn func(K) (V, error)) func(K) (V, error) {
	var mu sync.Mutex
	cache := make(map[K]V)
	if fn == nil {
		return func(K) (V, error) { var zero V; return zero, ErrNilFunction }
	}
	return func(key K) (V, error) {
		mu.Lock()
		value, ok := cache[key]
		mu.Unlock()
		if ok {
			return value, nil
		}
		value, err := fn(key)
		if err != nil {
			var zero V
			return zero, err
		}
		mu.Lock()
		if existing, exists := cache[key]; exists {
			value = existing
		} else {
			cache[key] = value
		}
		mu.Unlock()
		return value, nil
	}
}

// Bind fixes the first argument of a two-argument function and returns a
// single-argument closure.
func Bind[A, B, C any](fn func(A, B) C, first A) func(B) C {
	return func(second B) C { return fn(first, second) }
}

// Compose returns a closure that applies first and then second.
func Compose[A, B, C any](first func(A) B, second func(B) C) func(A) C {
	return func(value A) C { return second(first(value)) }
}

// Counter returns a stateful closure that increments and returns its value.
func Counter(start int) func() int {
	value := start
	return func() int { value++; return value }
}

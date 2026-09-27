package functional

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestMemoizeCachesSuccessfulResults(t *testing.T) {
	var calls atomic.Int32
	cached := Memoize(func(key string) (string, error) {
		calls.Add(1)
		return key + "!", nil
	})
	for i := 0; i < 2; i++ {
		got, err := cached("go")
		if err != nil || got != "go!" {
			t.Fatalf("got %q, %v", got, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestMemoizeDoesNotCacheErrors(t *testing.T) {
	var calls atomic.Int32
	cached := Memoize(func(int) (string, error) { calls.Add(1); return "", errors.New("fail") })
	_, _ = cached(1)
	_, _ = cached(1)
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
}

func TestBindComposeAndCounter(t *testing.T) {
	add := Bind(func(a, b int) int { return a + b }, 5)
	if add(3) != 8 {
		t.Fatal("Bind returned wrong result")
	}
	pipeline := Compose(func(v int) int { return v * 2 }, func(v int) string { return string(rune('0' + v)) })
	if pipeline(2) != "4" {
		t.Fatal("Compose returned wrong result")
	}
	next := Counter(0)
	if next() != 1 || next() != 2 {
		t.Fatal("Counter did not retain state")
	}
}

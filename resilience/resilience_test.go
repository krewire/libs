package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryStopsAfterSuccess(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), RetryOptions{MaxAttempts: 3}, func(context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil || attempts != 2 {
		t.Fatalf("err=%v attempts=%d", err, attempts)
	}
}

func TestRetryHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Retry(ctx, RetryOptions{MaxAttempts: 3}, func(context.Context) error { return errors.New("fail") }); !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestCircuitBreakerOpensAndRecovers(t *testing.T) {
	breaker := NewCircuitBreaker(2, time.Millisecond)
	failure := errors.New("failure")
	for i := 0; i < 2; i++ {
		if err := breaker.Do(context.Background(), func(context.Context) error { return failure }); !errors.Is(err, failure) {
			t.Fatal(err)
		}
	}
	if err := breaker.Do(context.Background(), func(context.Context) error { return nil }); !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("err=%v", err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := breaker.Do(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if breaker.State() != Closed {
		t.Fatalf("state=%v", breaker.State())
	}
}

func TestNilOperations(t *testing.T) {
	if !errors.Is(Retry(context.Background(), RetryOptions{}, nil), ErrNilOperation) {
		t.Fatal("Retry did not reject nil operation")
	}
	if !errors.Is(NewCircuitBreaker(1, time.Second).Do(context.Background(), nil), ErrNilOperation) {
		t.Fatal("breaker did not reject nil operation")
	}
}

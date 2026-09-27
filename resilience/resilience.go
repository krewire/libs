// Package resilience provides failure-handling utilities for closure-based workloads.
package resilience

import (
	"context"
	"errors"
	"math/rand/v2"
	"sync"
	"time"
)

var (
	ErrNilOperation = errors.New("resilience: nil operation")
	ErrCircuitOpen  = errors.New("resilience: circuit breaker is open")
)

type RetryOptions struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64
	ShouldRetry  func(error) bool
}

func (o RetryOptions) delay(attempt int) time.Duration {
	if o.InitialDelay <= 0 {
		return 0
	}
	multiplier := o.Multiplier
	if multiplier < 1 {
		multiplier = 2
	}
	delay := float64(o.InitialDelay)
	for i := 1; i < attempt; i++ {
		delay *= multiplier
	}
	if o.MaxDelay > 0 && time.Duration(delay) > o.MaxDelay {
		delay = float64(o.MaxDelay)
	}
	if o.Jitter > 0 {
		delay *= 1 + (rand.Float64()*2-1)*o.Jitter
	}
	if delay < 0 {
		return 0
	}
	return time.Duration(delay)
}

// Retry executes operation until it succeeds, the policy is exhausted, or ctx is cancelled.
func Retry(ctx context.Context, options RetryOptions, operation func(context.Context) error) error {
	if operation == nil {
		return ErrNilOperation
	}
	if ctx == nil {
		ctx = context.Background()
	}
	attempts := options.MaxAttempts
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := operation(ctx)
		if err == nil {
			return nil
		}
		last = err
		if options.ShouldRetry != nil && !options.ShouldRetry(err) {
			return err
		}
		if attempt == attempts {
			break
		}
		timer := time.NewTimer(options.delay(attempt))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return last
}

type BreakerState uint8

const (
	Closed BreakerState = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	mu               sync.Mutex
	failureThreshold int
	cooldown         time.Duration
	failures         int
	state            BreakerState
	openedAt         time.Time
	probeInFlight    bool
}

func NewCircuitBreaker(failureThreshold int, cooldown time.Duration) *CircuitBreaker {
	if failureThreshold < 1 {
		failureThreshold = 1
	}
	if cooldown < 0 {
		cooldown = 0
	}
	return &CircuitBreaker{failureThreshold: failureThreshold, cooldown: cooldown}
}

func (b *CircuitBreaker) State() BreakerState {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.transitionLocked(time.Now())
	return b.state
}

// Do executes operation when the circuit allows it and updates breaker state from its result.
func (b *CircuitBreaker) Do(ctx context.Context, operation func(context.Context) error) error {
	if b == nil {
		return ErrCircuitOpen
	}
	if operation == nil {
		return ErrNilOperation
	}
	if ctx == nil {
		ctx = context.Background()
	}
	b.mu.Lock()
	b.transitionLocked(time.Now())
	if b.state == Open || (b.state == HalfOpen && b.probeInFlight) {
		b.mu.Unlock()
		return ErrCircuitOpen
	}
	if b.state == HalfOpen {
		b.probeInFlight = true
	}
	b.mu.Unlock()
	err := operation(ctx)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.probeInFlight = false
	if err == nil {
		b.failures = 0
		b.state = Closed
		return nil
	}
	b.failures++
	if b.failures >= b.failureThreshold {
		b.state = Open
		b.openedAt = time.Now()
	}
	return err
}

func (b *CircuitBreaker) transitionLocked(now time.Time) {
	if b.state == Open && now.Sub(b.openedAt) >= b.cooldown {
		b.state = HalfOpen
		b.failures = 0
	}
}

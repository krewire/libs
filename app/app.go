// Package app provides application bootstrap, lifecycle, and service
// container primitives.
package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/krewire/libs/service"
)

// Container stores named application services.
type Container struct {
	mu       sync.RWMutex
	services map[string]any
}

// NewContainer creates an empty service container.
func NewContainer() *Container { return &Container{services: make(map[string]any)} }

// Set registers a service under name. Duplicate names are rejected.
func (c *Container) Set(name string, value any) error {
	if c == nil {
		return fmt.Errorf("app: nil container")
	}
	if name == "" {
		return fmt.Errorf("app: service name is required")
	}
	if value == nil {
		return fmt.Errorf("app: service %q must not be nil", name)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.services[name]; exists {
		return fmt.Errorf("app: service %q already registered", name)
	}
	c.services[name] = value
	return nil
}

// Get resolves a service by name.
func (c *Container) Get(name string) (any, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.services[name]
	return value, ok
}

// Resolve returns a typed service by name.
func Resolve[T any](c *Container, name string) (T, bool) {
	var zero T
	value, ok := c.Get(name)
	if !ok {
		return zero, false
	}
	resolved, ok := value.(T)
	return resolved, ok
}

// Runner executes the application-specific runtime after providers are ready.
// The application does not assume whether the runtime is HTTP, CLI, worker,
// desktop, or another execution model.
type Runner interface {
	Run(context.Context, *Container) error
}

// RunnerFunc adapts a function into a Runner.
type RunnerFunc func(context.Context, *Container) error

// Run implements Runner.
func (f RunnerFunc) Run(ctx context.Context, container *Container) error {
	if f == nil {
		return fmt.Errorf("app: nil runner")
	}
	return f(ctx, container)
}

// Application manages provider registration and lifecycle.
type Application struct {
	container *Container
	providers []service.Provider
	started   []service.Provider
	mu        sync.Mutex
}

// New creates an application with an empty service container.
func New() *Application { return &Application{container: NewContainer()} }

// Container returns the application's service container.
func (a *Application) Container() *Container { return a.container }

// Use adds providers in startup order. Providers must have unique names.
func (a *Application) Use(providers ...service.Provider) error {
	if a == nil {
		return fmt.Errorf("app: nil application")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	seen := make(map[string]struct{}, len(a.providers))
	for _, existing := range a.providers {
		seen[existing.Name()] = struct{}{}
	}
	for _, provider := range providers {
		if provider == nil {
			return fmt.Errorf("app: provider must not be nil")
		}
		name := provider.Name()
		if name == "" {
			return fmt.Errorf("app: provider name is required")
		}
		if _, exists := seen[name]; exists {
			return fmt.Errorf("app: provider %q already registered", name)
		}
		seen[name] = struct{}{}
		a.providers = append(a.providers, provider)
	}
	return nil
}

// Bootstrap registers all providers, then starts providers that implement
// service.Starter. A failed startup attempts to stop providers already started.
func (a *Application) Bootstrap(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("app: nil application")
	}
	a.mu.Lock()
	providers := append([]service.Provider(nil), a.providers...)
	a.mu.Unlock()
	for _, provider := range providers {
		if err := provider.Register(a.container); err != nil {
			return fmt.Errorf("app: register %q: %w", provider.Name(), err)
		}
	}
	for _, provider := range providers {
		starter, ok := provider.(service.Starter)
		if !ok {
			a.mu.Lock()
			a.started = append(a.started, provider)
			a.mu.Unlock()
			continue
		}
		if err := starter.Start(ctx, a.container); err != nil {
			_ = a.stopStarted(context.Background())
			return fmt.Errorf("app: start %q: %w", provider.Name(), err)
		}
		a.mu.Lock()
		a.started = append(a.started, provider)
		a.mu.Unlock()
	}
	return nil
}

// Run bootstraps providers, executes runner, and shuts providers down when the
// runner returns. Runner and shutdown errors are both preserved.
func (a *Application) Run(ctx context.Context, runner Runner) error {
	if a == nil {
		return fmt.Errorf("app: nil application")
	}
	if runner == nil {
		return fmt.Errorf("app: runner must not be nil")
	}
	if err := a.Bootstrap(ctx); err != nil {
		return err
	}
	runErr := runner.Run(ctx, a.container)
	stopErr := a.stopStarted(ctx)
	return errors.Join(runErr, stopErr)
}

// Shutdown stops started providers in reverse startup order.
func (a *Application) Shutdown(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("app: nil application")
	}
	return a.stopStarted(ctx)
}

func (a *Application) stopStarted(ctx context.Context) error {
	a.mu.Lock()
	started := append([]service.Provider(nil), a.started...)
	a.started = nil
	a.mu.Unlock()
	var stopErr error
	for i := len(started) - 1; i >= 0; i-- {
		if stopper, ok := started[i].(service.Stopper); ok {
			if err := stopper.Stop(ctx, a.container); err != nil {
				stopErr = fmt.Errorf("app: stop %q: %w", started[i].Name(), err)
			}
		}
	}
	return stopErr
}

var _ service.Registry = (*Container)(nil)

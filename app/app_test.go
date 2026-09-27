package app

import (
	"context"
	"testing"

	"github.com/krewire/libs/service"
)

type provider struct {
	name             string
	started, stopped bool
}

func (p *provider) Name() string                                  { return p.name }
func (p *provider) Register(r service.Registry) error             { return r.Set(p.name, p) }
func (p *provider) Start(context.Context, service.Registry) error { p.started = true; return nil }
func (p *provider) Stop(context.Context, service.Registry) error  { p.stopped = true; return nil }

func TestBootstrapResolvesAndShutsDownProviders(t *testing.T) {
	a := New()
	p := &provider{name: "clock"}
	if err := a.Use(p); err != nil {
		t.Fatal(err)
	}
	if err := a.Bootstrap(context.Background()); err != nil {
		t.Fatal(err)
	}
	got, ok := Resolve[*provider](a.Container(), "clock")
	if !ok || got != p || !p.started {
		t.Fatalf("provider was not started: got=%v started=%v", got, p.started)
	}
	if err := a.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !p.stopped {
		t.Fatal("provider was not stopped")
	}
}

func TestDuplicateProviderRejected(t *testing.T) {
	a := New()
	p := &provider{name: "clock"}
	if err := a.Use(p, &provider{name: "clock"}); err == nil {
		t.Fatal("duplicate provider must fail")
	}
}

func TestCanceledBootstrapStillRegistersBeforeStart(t *testing.T) {
	a := New()
	p := &provider{name: "clock"}
	if err := a.Use(p); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := a.Bootstrap(ctx); err != nil {
		t.Fatalf("provider does not inspect context: %v", err)
	}
	if _, ok := a.Container().Get("clock"); !ok {
		t.Fatal("provider was not registered")
	}
}

func TestRunExecutesRunnerAndStopsProviders(t *testing.T) {
	a := New()
	p := &provider{name: "clock"}
	if err := a.Use(p); err != nil {
		t.Fatal(err)
	}
	called := false
	err := a.Run(context.Background(), RunnerFunc(func(_ context.Context, c *Container) error {
		called = true
		if _, ok := c.Get("clock"); !ok {
			t.Fatal("provider missing from runner")
		}
		return nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !called || !p.stopped {
		t.Fatalf("called=%v stopped=%v", called, p.stopped)
	}
}

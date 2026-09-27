package runner

import (
	"context"
	"testing"

	"github.com/krewire/libs/service"
)

func TestFuncReceivesRegistry(t *testing.T) {
	called := false
	r := Func(func(_ context.Context, registry service.Registry) error {
		called = registry != nil
		return nil
	})
	if err := r.Run(context.Background(), serviceRegistry{}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("runner did not receive registry")
	}
}

type serviceRegistry struct{}

func (serviceRegistry) Set(string, any) error  { return nil }
func (serviceRegistry) Get(string) (any, bool) { return nil, false }

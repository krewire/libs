package kern

import (
	"context"
	"testing"

	"github.com/krewire/libs/core"
)

type testModule struct {
	name string
	init func(*Kernel) error
}

func (m testModule) Name() string { return m.name }
func (m testModule) Init(k *Kernel) error {
	if m.init != nil {
		return m.init(k)
	}
	return nil
}

func TestNewValidatesProject(t *testing.T) {
	_, err := New(core.Project{Name: "", Kind: core.KindApp})
	if err == nil {
		t.Error("New should fail for invalid project")
	}
	_, err = New(core.Project{Name: "demo", Kind: core.KindApp})
	if err != nil {
		t.Errorf("New error = %v, want nil", err)
	}
}

func TestRegistryDuplicate(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(testModule{name: "a"}); err != nil {
		t.Fatalf("Register error = %v", err)
	}
	if err := r.Register(testModule{name: "a"}); err == nil {
		t.Error("duplicate Register should fail")
	}
}

func TestKernelBootAndExecute(t *testing.T) {
	k, _ := New(core.Project{Name: "demo", Kind: core.KindApp})
	called := false
	k.Use(testModule{name: "mod", init: func(k *Kernel) error {
		k.RegisterHandler(core.KindApp, func(ctx context.Context, w core.Workload) core.ExitCode {
			called = true
			return core.ExitCodeSuccess
		})
		return nil
	}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatalf("Boot error = %v", err)
	}
	w, _ := core.WorkloadFor(core.KindApp)
	code := k.Execute(context.Background(), w)
	if code != core.ExitCodeSuccess || !called {
		t.Errorf("Execute = %v called=%v, want success true", code, called)
	}
	// Unknown workload
	unknown := core.Workload{Kind: "unknown"}
	if code := k.Execute(context.Background(), unknown); code != core.ExitCodeUsage {
		t.Errorf("Execute unknown = %v, want Usage", code)
	}
}

func TestSupervisor(t *testing.T) {
	s := NewSupervisor()
	hit := false
	s.OnReload(func() { hit = true })
	s.Reload()
	if !hit {
		t.Error("OnReload not triggered")
	}
	// Start/Stop with no modules should succeed
	if err := s.Start(context.Background()); err != nil {
		t.Errorf("Start error = %v", err)
	}
	if err := s.Stop(context.Background()); err != nil {
		t.Errorf("Stop error = %v", err)
	}
}

package core

import "testing"

func TestProjectValidate(t *testing.T) {
	cases := []struct {
		p  Project
		ok bool
	}{
		{Project{Name: "demo", Kind: KindApp}, true},
		{Project{Name: "my-app", Kind: KindCLI, ModulePath: "example.com/my-app"}, true},
		{Project{Name: "", Kind: KindApp}, false},
		{Project{Name: "Bad_Name", Kind: KindApp}, false},
		{Project{Name: "demo", Kind: Kind("unknown")}, false},
		{Project{Name: "demo", Kind: KindApp, ConfigPath: "krewire.yaml"}, true},
		{Project{Name: "demo", Kind: KindApp, ConfigPath: "ssg.yaml"}, false},
	}
	for i, c := range cases {
		err := c.p.Validate()
		if c.ok && err != nil {
			t.Errorf("case %d Validate error = %v, want nil", i, err)
		}
		if !c.ok && err == nil {
			t.Errorf("case %d Validate succeeded, want error", i)
		}
	}
}

func TestValidateKrewireYamlPath(t *testing.T) {
	if err := ValidateKrewireYamlPath("krewire.yaml"); err != nil {
		t.Errorf("valid path error = %v", err)
	}
	if err := ValidateKrewireYamlPath("./krewire.yaml"); err != nil {
		t.Errorf("valid path with ./ error = %v", err)
	}
	if err := ValidateKrewireYamlPath("ssg.yaml"); err == nil {
		t.Error("ssg.yaml should fail")
	}
}

func TestIsOptIn(t *testing.T) {
	if !IsOptIn(KindApp, []string{"github.com/krewire/framework/tui"}) {
		t.Error("app importing tui should be opt-in true")
	}
	if IsOptIn(KindApp, []string{"github.com/krewire/framework/service"}) {
		t.Error("app importing service should be opt-in false")
	}
	if !IsOptIn(KindService, []string{"github.com/krewire/framework/service"}) {
		t.Error("service importing service should be true")
	}
}

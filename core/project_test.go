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

func TestHasOptInViolation(t *testing.T) {
	cases := []struct {
		name     string
		kind     Kind
		imported []string
		want     bool
	}{
		{"app imports unrelated battery", KindApp, []string{"github.com/krewire/framework/tui"}, false},
		{"app imports service", KindApp, []string{"github.com/krewire/framework/service"}, true},
		{"app imports infra", KindApp, []string{"github.com/krewire/framework/infra"}, true},
		{"app imports service subpackage", KindApp, []string{"github.com/krewire/framework/service/gateway"}, true},
		{"app imports service among many", KindApp, []string{"github.com/krewire/framework/web", "github.com/krewire/framework/service"}, true},
		{"app imports worker in-process", KindApp, []string{"github.com/krewire/framework/worker"}, false},
		{"app imports runtime", KindApp, []string{"github.com/krewire/framework/runtime"}, false},
		{"prefix-only match is not a violation", KindApp, []string{"github.com/krewire/framework/serviceless"}, false},
		{"service kind imports service", KindService, []string{"github.com/krewire/framework/service"}, false},
		{"infra kind imports infra", KindInfra, []string{"github.com/krewire/framework/infra"}, false},
		{"site kind imports service", KindSite, []string{"github.com/krewire/framework/service"}, false},
		{"app imports nothing", KindApp, nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := HasOptInViolation(c.kind, c.imported); got != c.want {
				t.Errorf("HasOptInViolation(%q, %v) = %v, want %v", c.kind, c.imported, got, c.want)
			}
		})
	}
}

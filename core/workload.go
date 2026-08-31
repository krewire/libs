// Package core — workload registry (KWF-M8K2Q, KWL-K1N2Q).
package core

import (
	"fmt"
	"strings"
)

// Kind is a Krewire project kind. Eight kinds cover the unified workload spectrum.
type Kind string

const (
	KindApp     Kind = "app"
	KindCLI     Kind = "cli"
	KindSite    Kind = "site"
	KindBook    Kind = "book"
	KindWorker  Kind = "worker"
	KindService Kind = "service"
	KindInfra   Kind = "infra"
	KindKernel  Kind = "kernel"
)

// AllKinds lists every valid Kind in canonical order.
var AllKinds = []Kind{KindApp, KindCLI, KindSite, KindBook, KindWorker, KindService, KindInfra, KindKernel}

// IsValid reports whether k is one of the eight known kinds.
func (k Kind) IsValid() bool {
	switch k {
	case KindApp, KindCLI, KindSite, KindBook, KindWorker, KindService, KindInfra, KindKernel:
		return true
	default:
		return false
	}
}

// ParseKind parses s as a Kind, returning UsageError on unknown.
func ParseKind(s string) (Kind, error) {
	k := Kind(strings.TrimSpace(s))
	if !k.IsValid() {
		return "", UsageError(fmt.Sprintf("unknown project kind %q: want one of %s", s, strings.Join(kindsAsStrings(), ", ")))
	}
	return k, nil
}

func kindsAsStrings() []string {
	out := make([]string, len(AllKinds))
	for i, k := range AllKinds {
		out[i] = string(k)
	}
	return out
}

// Status is the implementation status of a workload.
type Status string

const (
	StatusShipped Status = "shipped"
	StatusPlanned Status = "planned"
)

// Workload describes one cell of the unified workload matrix.
type Workload struct {
	Kind    Kind   `json:"kind"`
	Package string `json:"package"`
	Title   string `json:"title"`
	SpecID  string `json:"specId"` // e.g. KWF-5XJFC
	Status  Status `json:"status"`
}

// Workloads is the canonical 9-workload matrix from internal/docs/project-vision.md.
var Workloads = []Workload{
	{Kind: KindCLI, Package: "framework/tui", Title: "CLI tools", SpecID: "KWF-5XJFC", Status: StatusShipped},
	{Kind: KindApp, Package: "framework/web", Title: "Backend / API", SpecID: "KWF-M07QS", Status: StatusShipped},
	{Kind: KindSite, Package: "framework/web/ssg", Title: "Static sites (SSG)", SpecID: "KWF-PT8OD", Status: StatusShipped},
	{Kind: KindBook, Package: "mdbind", Title: "Documentation sites", SpecID: "KWM-FX9H2", Status: StatusShipped},
	{Kind: KindApp, Package: "framework/app", Title: "Fullstack / Monolith", SpecID: "KWF-C4087", Status: StatusShipped},
	{Kind: KindSite, Package: "framework/runtime", Title: "Frontend (client)", SpecID: "KWF-T4X9P", Status: StatusShipped},
	{Kind: KindWorker, Package: "framework/worker", Title: "Workers & jobs", SpecID: "KWF-L5H2F", Status: StatusShipped},
	{Kind: KindService, Package: "framework/service", Title: "Microservice", SpecID: "KWF-L5H2F", Status: StatusShipped},
	{Kind: KindInfra, Package: "framework/infra", Title: "Cloud infrastructure", SpecID: "KWF-B7N3D", Status: StatusShipped},
}

// WorkloadFor returns the first Workload matching k and whether it was found.
func WorkloadFor(k Kind) (Workload, bool) {
	for _, w := range Workloads {
		if w.Kind == k {
			return w, true
		}
	}
	return Workload{}, false
}

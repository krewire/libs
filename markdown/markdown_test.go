// Tests for KWL-Q3N8P (KWL-MD-001..003) — Scope: Package libs/markdown.
package markdown

import (
	"strings"
	"testing"
)

func TestKWL_MD_001_Render_GFMAndHeadingIDs(t *testing.T) {
	src := []byte("# Hello\n\nA **bold** and ~~strike~~.\n\n| a | b |\n|---|---|\n| 1 | 2 |\n")
	html, err := Render(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<h1 id="hello">Hello</h1>`,
		`<strong>bold</strong>`,
		`<del>strike</del>`,
		"<table>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("Render output missing %q:\n%s", want, html)
		}
	}
}

func TestKWL_MD_002_RenderWithBase_PrefixesLinks(t *testing.T) {
	src := []byte("[next](/two) and ![alt](/img.png)")
	html, err := RenderWithBase(src, "/guide/")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `href="/guide/two"`) || !strings.Contains(html, `src="/guide/img.png"`) {
		t.Errorf("RenderWithBase = %s", html)
	}
	root, err := RenderWithBase([]byte("[x](/y)"), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(root, `href="/y"`) {
		t.Errorf("empty base must not rewrite: %s", root)
	}
}

func TestKWL_MD_003_PrefixLinks(t *testing.T) {
	tt := []struct{ in, base, want string }{
		{`<a href="/getting-started">`, "/guide/", `<a href="/guide/getting-started">`},
		{`<a href="/">`, "/guide/", `<a href="/guide/">`},
		{`<a href="/guide/y">`, "/guide/", `<a href="/guide/y">`},
		{`<img src="/img.png">`, "/guide/", `<img src="/guide/img.png">`},
		{`<a href="/x">`, "/", `<a href="/x">`},
		{`<a href="//example.com">`, "/guide/", `<a href="//example.com">`},
	}
	for _, tc := range tt {
		if got := PrefixLinks(tc.in, tc.base); got != tc.want {
			t.Errorf("PrefixLinks(%q,%q) = %q, want %q", tc.in, tc.base, got, tc.want)
		}
	}
}

func TestKWL_MD_003_RenderDeterministic(t *testing.T) {
	src := []byte("# T\n\nBody.\n")
	a, _ := Render(src)
	b, _ := Render(src)
	if a != b {
		t.Error("identical input must render identically")
	}
}

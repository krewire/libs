// Package markdown provides the shared Goldmark-based Markdown renderer for
// the Krewire ecosystem. Both mdbind (book) and framework/web/ssg + dsl use it
// so a docs site can start as a lightweight manuscript and progressively
// enhance to a full ssg site without re-parsing or duplicated dependencies.
// It is the leaf Package that allows mdbind and framework to be depended on
// together (KWF-M8K2Q, KWM-FX9H2).
package markdown

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var gold = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

var absLinkRe = regexp.MustCompile(`(href|src)="/([^"]*)"`)

// Render converts Markdown src to HTML using Goldmark GFM + AutoHeadingID.
func Render(src []byte) (string, error) {
	var buf bytes.Buffer
	if err := gold.Convert(src, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderWithBase converts Markdown src to HTML and rewrites absolute
// href/src="/*" links to be resolved under base (e.g. "/guide/").
// Page links keep their extensionless form; only the site root keeps the
// trailing slash. When base is "/" or empty no rewriting occurs.
func RenderWithBase(src []byte, base string) (string, error) {
	html, err := Render(src)
	if err != nil {
		return "", err
	}
	return prefixLinks(html, base), nil
}

// PrefixLinks rewrites absolute href/src links so they resolve under base.
// Exported for book tests and progressive callers; internal prefixLinks
// delegates to it.
func PrefixLinks(html, base string) string {
	return prefixLinks(html, base)
}

var codeBlockRe = regexp.MustCompile(`(?is)(<pre\b[^>]*>.*?</pre>|<code\b[^>]*>.*?</code>)`)

// prefixLinks rewrites absolute href/src links so they resolve under base,
// skipping <pre> and <code> blocks to prevent modifying code examples.
func prefixLinks(html, base string) string {
	prefix := ""
	if base != "" && base != "/" {
		prefix = strings.Trim(base, "/")
	}
	if prefix == "" {
		return html
	}

	matches := codeBlockRe.FindAllStringIndex(html, -1)
	if len(matches) == 0 {
		return rewriteLinks(html, prefix)
	}

	var sb strings.Builder
	lastIdx := 0
	for _, m := range matches {
		if m[0] > lastIdx {
			sb.WriteString(rewriteLinks(html[lastIdx:m[0]], prefix))
		}
		sb.WriteString(html[m[0]:m[1]])
		lastIdx = m[1]
	}
	if lastIdx < len(html) {
		sb.WriteString(rewriteLinks(html[lastIdx:], prefix))
	}
	return sb.String()
}

func rewriteLinks(html, prefix string) string {
	return absLinkRe.ReplaceAllStringFunc(html, func(m string) string {
		sub := absLinkRe.FindStringSubmatch(m)
		attr, rest := sub[1], sub[2]
		if strings.HasPrefix(rest, "/") || (prefix != "" && (strings.HasPrefix(rest, prefix+"/") || rest == prefix)) {
			return m
		}
		out := attr + `="/`
		if prefix != "" {
			out += prefix + "/"
		}
		out += rest
		return out + `"`
	})
}

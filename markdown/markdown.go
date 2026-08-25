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

// prefixLinks rewrites absolute href/src links so they resolve under base.
func prefixLinks(html, base string) string {
	prefix := ""
	if base != "" && base != "/" {
		prefix = strings.Trim(base, "/")
	}
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

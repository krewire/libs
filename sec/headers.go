package sec

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// SecurityOptions tunes the security-headers middleware.
type SecurityOptions struct {
	CSP               string
	Frame             string
	HSTS              int
	PermissionsPolicy string
}

// StripTags removes HTML tags from s using an HTML tokenizer in an iterative
// loop to prevent nested tag evasion (e.g. <<script>script>) and quoted attribute
// bypasses. Escaping remains html/template's job.
func StripTags(s string) string {
	prev := s
	for i := 0; i < 10; i++ {
		stripped := stripHTMLOnce(prev)
		if stripped == prev {
			break
		}
		prev = stripped
	}
	return strings.TrimSpace(prev)
}

func stripHTMLOnce(s string) string {
	var sb strings.Builder
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return sb.String()
		case html.TextToken:
			sb.Write(z.Text())
		}
	}
}

// SecurityHeaders returns middleware applying browser hardening headers.
func SecurityHeaders(opts ...func(*SecurityOptions)) Middleware {
	o := &SecurityOptions{}
	for _, f := range opts {
		f(o)
	}
	csp := o.CSP
	if csp == "" {
		csp = "default-src 'self'"
	}
	frame := o.Frame
	if frame == "" {
		frame = "DENY"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			setOnce(h, "X-Content-Type-Options", "nosniff")
			setOnce(h, "X-Frame-Options", frame)
			setOnce(h, "Referrer-Policy", "strict-origin-when-cross-origin")
			setOnce(h, "Content-Security-Policy", csp)
			if o.HSTS > 0 {
				setOnce(h, "Strict-Transport-Security", "max-age="+strconvItoa(o.HSTS)+"; includeSubDomains")
			}
			if o.PermissionsPolicy != "" {
				setOnce(h, "Permissions-Policy", o.PermissionsPolicy)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func setOnce(h http.Header, key, val string) {
	if h.Get(key) == "" {
		h.Set(key, val)
	}
}

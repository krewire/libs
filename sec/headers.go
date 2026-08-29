package sec

import (
	"net/http"
	"regexp"
	"strings"
)

// SecurityOptions tunes the security-headers middleware.
type SecurityOptions struct {
	CSP               string
	Frame             string
	HSTS              int
	PermissionsPolicy string
}

var tagStripper = regexp.MustCompile(`<[^>]*>`)

// StripTags removes HTML tags from s — a defense-in-depth helper for fields
// that must be plain text. Escaping remains html/template's job.
func StripTags(s string) string {
	return strings.TrimSpace(tagStripper.ReplaceAllString(s, ""))
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

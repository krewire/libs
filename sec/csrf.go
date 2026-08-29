package sec

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
)

// CSRFOptions tunes CSRF protection.
type CSRFOptions struct {
	CookieName string
	HeaderName string
	FieldName  string
	Secure     bool
	HTTPOnly   bool
}

type csrfCtxKey struct{}

// CSRF returns double-submit token middleware.
func CSRF(opts ...func(*CSRFOptions)) Middleware {
	o := &CSRFOptions{}
	for _, f := range opts {
		f(o)
	}
	if o.CookieName == "" {
		o.CookieName = "XSRF-TOKEN"
	}
	if o.HeaderName == "" {
		o.HeaderName = "X-CSRF-Token"
	}
	if o.FieldName == "" {
		o.FieldName = "csrf_token"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := cookieVal(r, o.CookieName)
			safe := r.Method == http.MethodGet || r.Method == http.MethodHead ||
				r.Method == http.MethodOptions || r.Method == http.MethodTrace

			if token == "" && safe && r.Header.Get("Authorization") == "" {
				token = randomToken()
				http.SetCookie(w, &http.Cookie{
					Name:     o.CookieName,
					Value:    token,
					Path:     "/",
					Secure:   o.Secure,
					HttpOnly: o.HTTPOnly,
					SameSite: http.SameSiteLaxMode,
				})
			}

			if !safe && r.Header.Get("Authorization") == "" {
				submitted := r.Header.Get(o.HeaderName)
				if submitted == "" {
					_ = r.ParseForm()
					submitted = r.PostForm.Get(o.FieldName)
				}
				if submitted == "" || token == "" || !constantTimeEqual(submitted, token) {
					Error(w, Forbidden("csrf token mismatch"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), csrfCtxKey{}, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CSRFFrom returns the current request's CSRF token.
func CSRFFrom(ctx context.Context) string {
	v, _ := ctx.Value(csrfCtxKey{}).(string)
	return v
}

func constantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("sec: crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

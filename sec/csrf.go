package sec

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

// CSRFOptions tunes CSRF protection.
type CSRFOptions struct {
	CookieName          string
	HeaderName          string
	FieldName           string
	Secure              bool
	HTTPOnly            bool
	PrefixHost          bool   // If true, prefixes cookie with __Host- and forces Secure + Path=/
	Secret              []byte // Optional HMAC secret to sign cookie tokens against cookie tossing
	IgnoreAuthorization bool   // If true, requests with Authorization header skip CSRF check
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
	if o.PrefixHost {
		if !strings.HasPrefix(o.CookieName, "__Host-") {
			o.CookieName = "__Host-" + o.CookieName
		}
		o.Secure = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieValStr := cookieVal(r, o.CookieName)
			safe := r.Method == http.MethodGet || r.Method == http.MethodHead ||
				r.Method == http.MethodOptions || r.Method == http.MethodTrace

			rawToken := cookieValStr
			if len(o.Secret) > 0 && cookieValStr != "" {
				parts := strings.Split(cookieValStr, ".")
				if len(parts) == 2 && validHMAC(parts[0], parts[1], o.Secret) {
					rawToken = parts[0]
				} else {
					rawToken = ""
				}
			}

			if rawToken == "" && safe {
				rawToken = randomToken()
				cookieValToSet := rawToken
				if len(o.Secret) > 0 {
					cookieValToSet = rawToken + "." + signHMAC(rawToken, o.Secret)
				}
				http.SetCookie(w, &http.Cookie{
					Name:     o.CookieName,
					Value:    cookieValToSet,
					Path:     "/",
					Secure:   o.Secure,
					HttpOnly: o.HTTPOnly,
					SameSite: http.SameSiteLaxMode,
				})
			}

			skip := o.IgnoreAuthorization && r.Header.Get("Authorization") != ""
			if !safe && !skip {
				submitted := r.Header.Get(o.HeaderName)
				if submitted == "" {
					_ = r.ParseForm()
					submitted = r.PostForm.Get(o.FieldName)
				}
				if submitted == "" || rawToken == "" || !constantTimeEqual(submitted, rawToken) {
					Error(w, Forbidden("csrf token mismatch"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), csrfCtxKey{}, rawToken)
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

func signHMAC(token string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func validHMAC(token, sig string, secret []byte) bool {
	expected := signHMAC(token, secret)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(sig)) == 1
}

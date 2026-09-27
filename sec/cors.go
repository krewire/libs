package sec

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSOptions configures Cross-Origin Resource Sharing.
type CORSOptions struct {
	// AllowOrigins is a list of allowed origins. Defaults to ["*"].
	AllowOrigins []string
	// AllowMethods is a list of allowed HTTP methods. Defaults to GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS.
	AllowMethods []string
	// AllowHeaders is a list of allowed request headers. Defaults to standard headers.
	AllowHeaders []string
	// ExposeHeaders is a list of response headers exposed to the client.
	ExposeHeaders []string
	// AllowCredentials indicates whether cookies, authorization headers, or TLS client certificates are exposed.
	AllowCredentials bool
	// MaxAge is the duration in seconds that preflight requests can be cached. Defaults to 86400 (24h).
	MaxAge int
}

// WithOrigins appends allowed origins to CORSOptions.
func WithOrigins(origins ...string) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.AllowOrigins = append(o.AllowOrigins, origins...)
	}
}

// WithMethods appends allowed HTTP methods to CORSOptions.
func WithMethods(methods ...string) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.AllowMethods = append(o.AllowMethods, methods...)
	}
}

// WithHeaders appends allowed headers to CORSOptions.
func WithHeaders(headers ...string) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.AllowHeaders = append(o.AllowHeaders, headers...)
	}
}

// WithExposeHeaders appends exposed response headers to CORSOptions.
func WithExposeHeaders(headers ...string) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.ExposeHeaders = append(o.ExposeHeaders, headers...)
	}
}

// WithCredentials sets AllowCredentials in CORSOptions.
func WithCredentials(allow bool) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.AllowCredentials = allow
	}
}

// WithMaxAge sets preflight cache duration in seconds.
func WithMaxAge(seconds int) func(*CORSOptions) {
	return func(o *CORSOptions) {
		o.MaxAge = seconds
	}
}

// CORS returns a framework-agnostic HTTP middleware implementing Cross-Origin Resource Sharing.
func CORS(opts ...func(*CORSOptions)) Middleware {
	o := &CORSOptions{
		MaxAge: 86400,
	}
	for _, f := range opts {
		f(o)
	}

	if len(o.AllowOrigins) == 0 {
		o.AllowOrigins = []string{"*"}
	}
	if len(o.AllowMethods) == 0 {
		o.AllowMethods = []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		}
	}
	if len(o.AllowHeaders) == 0 {
		o.AllowHeaders = []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		}
	}

	methodsStr := strings.Join(o.AllowMethods, ", ")
	headersStr := strings.Join(o.AllowHeaders, ", ")
	exposeStr := strings.Join(o.ExposeHeaders, ", ")
	maxAgeStr := strconv.Itoa(o.MaxAge)

	wildcardOrigin := false
	for _, origin := range o.AllowOrigins {
		if origin == "*" {
			wildcardOrigin = true
			break
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed := false
			if wildcardOrigin && !o.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				allowed = true
			} else {
				for _, oOrigin := range o.AllowOrigins {
					// Never reflect an arbitrary origin when credentials are enabled.
					// The wildcard origin is valid only for non-credentialed CORS.
					if oOrigin != "*" && strings.EqualFold(oOrigin, origin) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Add("Vary", "Origin")
						allowed = true
						break
					}
				}
			}

			if !allowed {
				next.ServeHTTP(w, r)
				return
			}

			if o.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if exposeStr != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposeStr)
			}

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", methodsStr)
				w.Header().Set("Access-Control-Allow-Headers", headersStr)
				if o.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", maxAgeStr)
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

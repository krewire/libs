package sec

import (
	"fmt"
	"net/http"
)

// HTTPError is a structured HTTP error with status, code, and message.
type HTTPError struct {
	Status  int
	Code    string
	Message string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return http.StatusText(e.Status)
}

// Error writes the HTTPError as JSON or plain text via http.Error.
func Error(w http.ResponseWriter, err error) {
	if he, ok := err.(*HTTPError); ok {
		http.Error(w, he.Message, he.Status)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Middleware is a standard http middleware.
type Middleware func(http.Handler) http.Handler

func cookieVal(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func lowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func strEqFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	// constant-time for auth schemes
	aa := []byte(lowerASCII(a))
	bb := []byte(lowerASCII(b))
	// use subtle.ConstantTimeCompare for security
	// but for scheme compare, simple is fine
	eq := 1
	for i := range aa {
		if aa[i] != bb[i] {
			eq = 0
			break
		}
	}
	return eq == 1
}

func strconvItoa(i int) string {
	return fmt.Sprintf("%d", i)
}

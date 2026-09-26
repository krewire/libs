package sec

import (
	"net/http"

	"github.com/krewire/libs/auth"
)

// Error writes the HTTPError as JSON or plain text via http.Error.
func Error(w http.ResponseWriter, err error) {
	auth.Error(w, err)
}

// Middleware is a standard http middleware.
type Middleware = auth.Middleware

// cookieVal is re-exported from auth for backward compatibility within sec.
func cookieVal(r *http.Request, name string) string {
	return auth.CookieValue(r, name)
}

// lowerASCII is re-exported from auth for backward compatibility within sec.
func lowerASCII(s string) string {
	return auth.LowerASCII(s)
}

// strEqFold is re-exported from auth for backward compatibility within sec.
func strEqFold(a, b string) bool {
	return auth.StrEqFold(a, b)
}

// strconvItoa is re-exported from auth for backward compatibility within sec.
func strconvItoa(i int) string {
	return auth.StrconvItoa(i)
}

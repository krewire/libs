package sec

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCSRFSafeMethods(t *testing.T) {
	mw := CSRF()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := CSRFFrom(r.Context())
		w.Write([]byte(tok))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/form", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	cookie := rec.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "XSRF-TOKEN=") {
		t.Errorf("Set-Cookie missing XSRF-TOKEN: %s", cookie)
	}
}

func TestCSRFUnsafeBlockedWithoutToken(t *testing.T) {
	mw := CSRF()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", rec.Code)
	}
}

func TestCSRFUnsafeAllowedWithValidHeader(t *testing.T) {
	mw := CSRF()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := "valid-csrf-token-1234567890abcdef"
	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	req.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: token})
	req.Header.Set("X-CSRF-Token", token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

package auth

import (
	"strings"
	"testing"
	"time"
)

func TestJWTRoundTrip(t *testing.T) {
	secret := []byte("secret-key")
	claims := DefaultClaims("user-123", time.Hour)
	claims["email"] = "user@example.com"

	token, err := SignJWT(secret, claims)
	if err != nil {
		t.Fatalf("SignJWT error = %v", err)
	}

	parsed, err := ParseJWT(secret, token)
	if err != nil {
		t.Fatalf("ParseJWT error = %v", err)
	}
	if parsed["sub"] != "user-123" {
		t.Errorf("sub = %v, want user-123", parsed["sub"])
	}
	if parsed["email"] != "user@example.com" {
		t.Errorf("email = %v, want user@example.com", parsed["email"])
	}
}

func TestJWTExpired(t *testing.T) {
	secret := []byte("secret-key")
	claims := DefaultClaims("user-123", -time.Minute)

	token, err := SignJWT(secret, claims)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ParseJWT(secret, token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestJWTTampered(t *testing.T) {
	secret := []byte("secret-key")
	claims := DefaultClaims("user-123", time.Hour)
	token, _ := SignJWT(secret, claims)

	parts := strings.Split(token, ".")
	tampered := parts[0] + "." + parts[1] + "extra." + parts[2]
	_, err := ParseJWT(secret, tampered)
	if err == nil {
		t.Error("expected error for tampered token")
	}
}

func TestB64JSON(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	encoded, err := B64JSON(data)
	if err != nil {
		t.Fatalf("B64JSON error = %v", err)
	}
	if encoded == "" {
		t.Error("B64JSON returned empty string")
	}
}

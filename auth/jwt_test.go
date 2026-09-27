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

func TestJWT_NBF_NotYetValid(t *testing.T) {
	secret := []byte("secret-key")
	claims := DefaultClaims("user-123", time.Hour)
	claims["nbf"] = time.Now().Add(10 * time.Minute).Unix()

	token, err := SignJWT(secret, claims)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ParseJWT(secret, token)
	if err == nil {
		t.Error("expected error for token with future nbf")
	}

	// With leeway >= 10m, it should pass
	_, err = ParseJWT(secret, token, WithLeeway(15*time.Minute))
	if err != nil {
		t.Errorf("expected token to pass with sufficient leeway: %v", err)
	}
}

func TestJWT_IAT_Future(t *testing.T) {
	secret := []byte("secret-key")
	claims := DefaultClaims("user-123", time.Hour)
	claims["iat"] = time.Now().Add(10 * time.Minute).Unix()

	token, err := SignJWT(secret, claims)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ParseJWT(secret, token)
	if err == nil {
		t.Error("expected error for token with future iat")
	}
}

func TestJWT_RequireExp(t *testing.T) {
	secret := []byte("secret-key")
	claims := Claims{"sub": "user-123"} // no exp

	token, _ := SignJWT(secret, claims)
	_, err := ParseJWT(secret, token, WithRequireExp(true))
	if err == nil {
		t.Error("expected error when exp is required but missing")
	}

	_, err = ParseJWT(secret, token, WithRequireExp(false))
	if err != nil {
		t.Errorf("expected success when exp is optional: %v", err)
	}
}

func TestJWT_MinSecretLength(t *testing.T) {
	shortSecret := []byte("short")
	claims := DefaultClaims("user-123", time.Hour)
	token, _ := SignJWT(shortSecret, claims)

	_, err := ParseJWT(shortSecret, token, WithMinSecretLength(32))
	if err == nil {
		t.Error("expected error for short secret when min length is 32")
	}

	validSecret := []byte("01234567890123456789012345678901") // 32 bytes
	tokenValid, _ := SignJWT(validSecret, claims)
	_, err = ParseJWT(validSecret, tokenValid, WithMinSecretLength(32))
	if err != nil {
		t.Errorf("expected success for 32-byte secret: %v", err)
	}
}

package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key-1234567890"
	username := "admin"
	expiry := 24 * time.Hour

	token, err := GenerateToken(username, secret, expiry)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	gotUsername, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if gotUsername != username {
		t.Errorf("ValidateToken() username = %q, want %q", gotUsername, username)
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	secret := "correct-secret"
	wrongSecret := "wrong-secret"

	token, _ := GenerateToken("admin", secret, 24*time.Hour)

	_, err := ValidateToken(token, wrongSecret)
	if err == nil {
		t.Fatal("ValidateToken() should fail with wrong secret")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	secret := "test-secret"

	token, _ := GenerateToken("admin", secret, -1*time.Hour)

	_, err := ValidateToken(token, secret)
	if err == nil {
		t.Fatal("ValidateToken() should fail for expired token")
	}
}

func TestValidateTokenInvalidFormat(t *testing.T) {
	secret := "test-secret"

	_, err := ValidateToken("not-a-valid-token", secret)
	if err == nil {
		t.Fatal("ValidateToken() should fail for invalid token format")
	}
}

func TestValidateTokenEmpty(t *testing.T) {
	secret := "test-secret"

	_, err := ValidateToken("", secret)
	if err == nil {
		t.Fatal("ValidateToken() should fail for empty token")
	}
}

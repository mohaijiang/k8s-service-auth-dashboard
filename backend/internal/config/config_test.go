package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.Namespace != "dashboard-auth-system" {
		t.Errorf("Namespace = %q, want %q", cfg.Namespace, "dashboard-auth-system")
	}
	if cfg.InitAdminUsername != "admin" {
		t.Errorf("InitAdminUsername = %q, want %q", cfg.InitAdminUsername, "admin")
	}
	if cfg.InitAdminPassword != "" {
		t.Errorf("InitAdminPassword should be empty by default, got %q", cfg.InitAdminPassword)
	}
	if cfg.JWTExpiry != 24*time.Hour {
		t.Errorf("JWTExpiry = %v, want %v", cfg.JWTExpiry, 24*time.Hour)
	}
	if cfg.CORSAllowOrigin != "http://localhost:3000" {
		t.Errorf("CORSAllowOrigin = %q, want %q", cfg.CORSAllowOrigin, "http://localhost:3000")
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "9090")
	os.Setenv("NAMESPACE", "custom-ns")
	os.Setenv("INIT_ADMIN_USERNAME", "superadmin")
	os.Setenv("INIT_ADMIN_PASSWORD", "secret123")
	os.Setenv("JWT_EXPIRY", "48h")
	os.Setenv("CORS_ALLOW_ORIGIN", "http://frontend:3000")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.Namespace != "custom-ns" {
		t.Errorf("Namespace = %q, want %q", cfg.Namespace, "custom-ns")
	}
	if cfg.InitAdminUsername != "superadmin" {
		t.Errorf("InitAdminUsername = %q, want %q", cfg.InitAdminUsername, "superadmin")
	}
	if cfg.InitAdminPassword != "secret123" {
		t.Errorf("InitAdminPassword = %q, want %q", cfg.InitAdminPassword, "secret123")
	}
	if cfg.JWTExpiry != 48*time.Hour {
		t.Errorf("JWTExpiry = %v, want %v", cfg.JWTExpiry, 48*time.Hour)
	}
	if cfg.CORSAllowOrigin != "http://frontend:3000" {
		t.Errorf("CORSAllowOrigin = %q, want %q", cfg.CORSAllowOrigin, "http://frontend:3000")
	}
}

func TestLoadInvalidJWTExpiry(t *testing.T) {
	os.Clearenv()
	os.Setenv("JWT_EXPIRY", "not-a-duration")

	cfg := Load()

	if cfg.JWTExpiry != 24*time.Hour {
		t.Errorf("JWTExpiry with invalid input should fallback to 24h, got %v", cfg.JWTExpiry)
	}
}

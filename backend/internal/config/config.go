package config

import (
	"os"
	"time"
)

// Config holds application configuration from environment variables.
type Config struct {
	Port              string
	Env               string
	Namespace         string
	InitAdminUsername string
	InitAdminPassword string
	JWTExpiry         time.Duration
	CORSAllowOrigin   string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		Env:               getEnv("ENV", "development"),
		Namespace:         getEnv("NAMESPACE", "dashboard-auth-system"),
		InitAdminUsername: getEnv("INIT_ADMIN_USERNAME", "admin"),
		InitAdminPassword: getEnv("INIT_ADMIN_PASSWORD", ""),
		JWTExpiry:         parseDuration(getEnv("JWT_EXPIRY", "24h")),
		CORSAllowOrigin:   getEnv("CORS_ALLOW_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

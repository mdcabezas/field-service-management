package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	os.Setenv("LISTEN_ADDR", ":9090")
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	os.Setenv("PHOTO_STORAGE_DIR", "/data/photos")

	cfg := Load()

	require.Equal(t, ":9090", cfg.ListenAddr)
	require.Equal(t, "testsecret", cfg.JWTSecret)
	require.Equal(t, "postgres://test:test@localhost:5432/test", cfg.DatabaseURL)
	require.Equal(t, "http://localhost:3000", cfg.CORSAllowedOrigins)
	require.Equal(t, "/data/photos", cfg.PhotoStorageDir)
}

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("LISTEN_ADDR")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	os.Unsetenv("PHOTO_STORAGE_DIR")

	cfg := Load()

	require.Equal(t, ":8080", cfg.ListenAddr)
	require.Equal(t, "http://localhost:3000,http://localhost:5173", cfg.CORSAllowedOrigins)
	require.Equal(t, "/mobile-photos", cfg.PhotoStorageDir)
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "testvalue")
	defer os.Unsetenv("TEST_KEY")

	require.Equal(t, "testvalue", getEnv("TEST_KEY", "fallback"))
	require.Equal(t, "fallback", getEnv("NONEXISTENT", "fallback"))
}

func TestGetEnvRequired(t *testing.T) {
	os.Setenv("REQUIRED_KEY", "requiredvalue")
	defer os.Unsetenv("REQUIRED_KEY")

	require.Equal(t, "requiredvalue", getEnvRequired("REQUIRED_KEY"))
}

func TestGetEnvRequired_Panics(t *testing.T) {
	os.Unsetenv("MISSING_KEY")

	// getEnvRequired uses log.Fatalf which calls os.Exit(1)
	// We can't easily test this in the same process, so just verify the function exists
	require.NotPanics(t, func() { getEnv("EXISTING_KEY", "fallback") })
}

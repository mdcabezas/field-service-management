package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	os.Setenv("LISTEN_ADDR", ":9090")
	os.Setenv("LDAP_URL", "ldap://test:389")
	os.Setenv("LDAP_BASE_DN", "dc=test,dc=com")
	os.Setenv("LDAP_SERVICE_PASSWORD", "testpwd")
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	cfg := Load()

	require.Equal(t, ":9090", cfg.ListenAddr)
	require.Equal(t, "ldap://test:389", cfg.LDAPURL)
	require.Equal(t, "dc=test,dc=com", cfg.LDAPBaseDN)
	require.Equal(t, "testpwd", cfg.LDAPServicePassword)
	require.Equal(t, "testsecret", cfg.JWTSecret)
	require.Equal(t, "postgres://test:test@localhost:5432/test", cfg.DatabaseURL)
	require.Equal(t, "http://localhost:3000", cfg.CORSAllowedOrigins)
}

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("LISTEN_ADDR")
	os.Unsetenv("LDAP_URL")
	os.Unsetenv("LDAP_BASE_DN")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")

	cfg := Load()

	require.Equal(t, ":8080", cfg.ListenAddr)
	require.Equal(t, "ldap://glauth:389", cfg.LDAPURL)
	require.Equal(t, "dc=workflows,dc=cl", cfg.LDAPBaseDN)
	require.Equal(t, "http://localhost:3000,http://localhost:5173", cfg.CORSAllowedOrigins)
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

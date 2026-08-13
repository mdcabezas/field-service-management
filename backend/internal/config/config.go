package config

import (
	"log"
	"os"
)

type Config struct {
	ListenAddr          string
	LDAPURL             string
	LDAPBaseDN          string
	LDAPServicePassword string
	JWTSecret           string
	DatabaseURL         string
	CORSAllowedOrigins  string
}

func Load() *Config {
	return &Config{
		ListenAddr:          getEnv("LISTEN_ADDR", ":8080"),
		LDAPURL:             getEnv("LDAP_URL", "ldap://glauth:389"),
		LDAPBaseDN:          getEnv("LDAP_BASE_DN", "dc=workflows,dc=cl"),
		LDAPServicePassword: getEnvRequired("LDAP_SERVICE_PASSWORD"),
		JWTSecret:           getEnvRequired("JWT_SECRET"),
		DatabaseURL:         getEnvRequired("DATABASE_URL"),
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvRequired(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

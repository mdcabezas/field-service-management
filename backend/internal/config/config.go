package config

import (
	"log"
	"os"
)

type Config struct {
	ListenAddr         string
	JWTSecret          string
	DatabaseURL        string
	CORSAllowedOrigins string
	PhotoStorageDir    string
}

func Load() *Config {
	return &Config{
		ListenAddr:         getEnv("LISTEN_ADDR", ":8080"),
		JWTSecret:          getEnvRequired("JWT_SECRET"),
		DatabaseURL:        getEnvRequired("DATABASE_URL"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		PhotoStorageDir:    getEnv("PHOTO_STORAGE_DIR", "/mobile-photos"),
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

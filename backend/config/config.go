package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	DatabaseURL string
	Port        string
	FrontendURL string
	Env         string
}

// Load reads the .env file (if present) then maps env vars into a Config.
// Missing DATABASE_URL is fatal; other fields fall back to sensible defaults.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found — using environment variables")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnvOrDefault("PORT", "8080"),
		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		Env:         getEnvOrDefault("ENV", "development"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("config: DATABASE_URL is required but not set")
	}

	return cfg
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

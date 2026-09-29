package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBDatabase  string
	SSLMode     string
	APIHost     string
	APIPort     string
	CORSOrigins []string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "vbevents"),
		DBPassword:  getEnv("DB_PASSWORD", "vbevents"),
		DBDatabase:  getEnv("DB_DATABASE", "events"),
		SSLMode:     getEnv("SSL_MODE", "disable"),
		APIHost:     getEnv("API_HOST", "localhost"),
		APIPort:     getEnv("API_PORT", "8200"),
		CORSOrigins: getEnvSlice("CORS_ORIGINS", []string{"http://localhost:5173"}),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBDatabase string
	SSLMode    string
	APIHost    string
	APIPort    string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "vbevents"),
		DBPassword: getEnv("DB_PASSWORD", "vbevents"),
		DBDatabase: getEnv("DB_DATABASE", "events"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
		APIHost:    getEnv("API_HOST", "localhost"),
		APIPort:    getEnv("API_PORT", "8200"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

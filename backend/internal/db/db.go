package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return db, nil
}

func getDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "vb-events")
	password := getEnv("DB_PASSWORD", "vbevents")
	database := getEnv("DB_DATABASE", "events")
	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + database
	return dsn
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

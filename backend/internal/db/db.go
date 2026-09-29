package db

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/config"
)

func InitDB(cfg config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBDatabase

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return db, nil
}

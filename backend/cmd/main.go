package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/db"
)

func main() {
	cfg := config.Load()

	_, err := db.InitDB(cfg)
	if err != nil {
		panic("failed to connect to database")
	}

	e := echo.New()
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(cfg.APIHost + ":" + cfg.APIPort); err != nil {
		slog.Error("failed to start API", "error", err)
	} else {
		slog.Info("API started", "host", cfg.APIHost, "port", cfg.APIPort)
	}
	slog.Info("shutting down API")
}

package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"

	_ "github.com/williamf6894/VB-Events/docs"
	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/db"
	"github.com/williamf6894/VB-Events/internal/handlers"
	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
	"github.com/williamf6894/VB-Events/internal/services"
)

//	@title			VB-Events API
//	@version		1.0
//	@description	Events management API
//	@BasePath		/
const swaggerInitializer = `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/swagger/doc.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};`

func main() {
	cfg := config.Load()

	database, err := db.InitDB(cfg)
	if err != nil {
		panic("failed to connect to database")
	}

	if err := database.AutoMigrate(&models.Event{}, &models.Participant{}); err != nil {
		panic("failed to migrate database")
	}

	participantRepo := repository.NewParticipantRepository(database)
	participantService := services.NewParticipantService(participantRepo)
	participantHandler := handlers.NewParticipantHandler(participantService)

	e := echo.New()
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.Recover())

	e.POST("/participants", participantHandler.Create)
	e.GET("/participants", participantHandler.List)
	e.GET("/participants/:id", participantHandler.FindByID)
	e.PUT("/participants/:id", participantHandler.Update)
	e.DELETE("/participants/:id", participantHandler.Delete)

	e.GET("/swagger", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	e.GET("/swagger/doc.json", func(c *echo.Context) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read swagger doc")
		}
		return c.Blob(http.StatusOK, "application/json", []byte(doc))
	})
	e.GET("/swagger/swagger-initializer.js", func(c *echo.Context) error {
		return c.Blob(http.StatusOK, "text/javascript; charset=utf-8", []byte(swaggerInitializer))
	})
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger", http.FileServer(http.FS(swaggerFiles.FS)))))

	if err := e.Start(cfg.APIHost + ":" + cfg.APIPort); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}

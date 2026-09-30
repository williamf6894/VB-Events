package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"

	_ "github.com/williamf6894/VB-Events/docs"
	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/db"
	"github.com/williamf6894/VB-Events/internal/handlers"
	"github.com/williamf6894/VB-Events/internal/middleware"
	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
	"github.com/williamf6894/VB-Events/internal/services"
)

// @title			VB-Events API
// @version		1.0
// @description	Events management API
// @BasePath		/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
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
	// Config
	cfg := config.Load()

	if cfg.JWTSecret == "" {
		panic("JWT_SECRET must be set (put it in .env or the environment)")
	}

	// Database
	database, err := db.InitDB(cfg)
	if err != nil {
		panic("failed to connect to database")
	}

	if err := database.AutoMigrate(&models.Event{}, &models.Participant{}); err != nil {
		panic("failed to migrate database")
	}

	// Wiring up system
	participantRepo := repository.NewParticipantRepository(database)
	participantService := services.NewParticipantService(participantRepo)
	participantHandler := handlers.NewParticipantHandler(participantService)

	eventRepo := repository.NewEventRepository(database)
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	authService := services.NewAuthService(participantRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.Auth(authService, participantRepo)
	healthHandler := handlers.NewHealthHandler(database)

	// Middleware

	e := echo.New()
	e.Pre(echoMiddleware.RemoveTrailingSlash())
	// Depending on how much you want to log and how many logs you are sending
	// you may want to comment out this request logger.
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{

		// Skipping /healthz because its noisy
		Skipper: func(c *echo.Context) bool {
			return c.Request().URL.Path == "/healthz"
		},
		LogStatus:   true,
		LogURI:      true,
		HandleError: true,
		LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger := c.Logger()
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Duration("latency", v.Latency),
				)
				return nil
			}

			logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("error", v.Error.Error()),
			)
			return nil
		},
	}))
	e.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		XSSProtection:         "1, mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'; 'unsafe-inline'; 'unsafe-eval' ",
	}))
	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(40.0)))
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       300,
	}))
	e.Use(echoMiddleware.Recover())

	// Health
	e.GET("/healthz", healthHandler.Check)

	// API
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)
	e.GET("/auth/me", authHandler.Me, authMiddleware)

	e.GET("/participants", participantHandler.List)
	e.GET("/participants/:id", participantHandler.FindByID)
	e.PUT("/participants/:id", participantHandler.Update)
	e.DELETE("/participants/:id", participantHandler.Delete)

	e.POST("/events", eventHandler.Create)
	e.GET("/events", eventHandler.List)
	e.GET("/events/:id", eventHandler.FindByID)
	e.GET("/events/name/:name", eventHandler.FindByName)
	e.PUT("/events/:id", eventHandler.Update)
	e.DELETE("/events/:id", eventHandler.Delete)
	e.POST("/events/:id/participants", eventHandler.Join, authMiddleware)
	e.POST("/events/:id/participants/:participantId", eventHandler.Invite, authMiddleware)
	e.DELETE("/events/:id/participants", eventHandler.Leave, authMiddleware)

	// Swagger Documentation
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

	// Start Server
	slog.Info("Starting server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = echo.StartConfig{Address: ":" + cfg.APIPort}.Start(ctx, e)
	if err != nil {
		slog.Error("server exited", "error", err)
	}
}

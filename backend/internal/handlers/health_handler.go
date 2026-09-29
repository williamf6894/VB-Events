package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	sqlDB, _ := db.DB()
	return &HealthHandler{db: sqlDB}
}

type healthResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// CheckHealth godoc
// @Summary      Service health check
// @Description  Returns service status and verifies database connectivity
// @Tags         health
// @Produce      json
// @Success      200 {object} handlers.healthResponse
// @Failure      503 {object} handlers.healthResponse
// @Router       /healthz [get]
func (h *HealthHandler) Check(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{}
	healthy := true

	if err := h.db.PingContext(ctx); err != nil {
		healthy = false
		checks["database"] = "down: " + err.Error()
	} else {
		checks["database"] = "ok"
	}

	status := "ok"
	code := http.StatusOK
	if !healthy {
		status = "unavailable"
		code = http.StatusServiceUnavailable
	}

	return c.JSON(code, healthResponse{Status: status, Checks: checks})
}

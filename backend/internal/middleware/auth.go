package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
	"github.com/williamf6894/VB-Events/internal/services"
)

const ParticipantKey = "participant"

func Auth(authService *services.AuthService, participantRepo *repository.ParticipantRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}

			userID, err := authService.ValidateToken(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			participant, err := participantRepo.FindByID(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unknown user")
			}

			c.Set(ParticipantKey, participant)
			return next(c)
		}
	}
}

func ParticipantFrom(c *echo.Context) *models.Participant {
	participant, ok := c.Get(ParticipantKey).(*models.Participant)
	if !ok {
		return nil
	}
	return participant
}

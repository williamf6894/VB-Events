package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/williamf6894/VB-Events/internal/middleware"
	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type participantRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" example:"min 8 characters"`
}

// RegisterParticipant godoc
// @Summary      Register a new participant
// @Description  Creates a participant account with a name, unique email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body participantRequest true "Participant to register"
// @Success      201 {object} models.Participant
// @Failure      400 {object} object "Invalid body, or missing name/email/password (min 8 chars)"
// @Failure      409 {object} object "Email already in use"
// @Failure      500 {object} object "Internal error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *echo.Context) error {
	var req participantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	participant := &models.Participant{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.authService.Create(participant); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidParticipant):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrEmailTaken):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create participant")
		}
	}

	return c.JSON(http.StatusCreated, participant)
}

// LoginParticipant godoc
// @Summary      Log in with email and password
// @Description  Returns a signed JWT valid for 24 hours
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body services.LoginRequest true "Credentials"
// @Success      200 {object} services.TokenResponse
// @Failure      400 {object} object "Invalid request body"
// @Failure      401 {object} object "Invalid credentials"
// @Failure      500 {object} object "Internal error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *echo.Context) error {
	var req services.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	token, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to log in")
	}
	return c.JSON(http.StatusOK, token)
}

// Me godoc
// @Summary      Get the currently authenticated participant
// @Tags         auth
// @Produce      json
// @Success      200 {object} models.Participant
// @Failure      401 {object} object "Not authenticated"
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *echo.Context) error {
	participant := middleware.ParticipantFrom(c)
	if participant == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}
	return c.JSON(http.StatusOK, participant)
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"uuid"

	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/services"
)

type ParticipantHandler struct {
	service *services.ParticipantService
}

func NewParticipantHandler(service *services.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}

type participantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateParticipant godoc
// @Summary      Create a new participant
// @Description  Registers a participant with a name and unique email
// @Tags         participants
// @Accept       json
// @Produce      json
// @Param        request body participantRequest true "Participant to create"
// @Success      201 {object} models.Participant
// @Failure      400 {object} object "Invalid body, or missing name/email"
// @Failure      409 {object} object "Email already in use"
// @Failure      500 {object} object "Internal error"
// @Router       /participants [post]
func (h *ParticipantHandler) Create(c *echo.Context) error {
	var req participantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	participant := &models.Participant{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.service.Create(participant); err != nil {
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

// ListParticipants godoc
// @Summary      List all participants
// @Tags         participants
// @Produce      json
// @Success      200 {array} models.Participant
// @Failure      500 {object} object "Internal error"
// @Router       /participants [get]
func (h *ParticipantHandler) List(c *echo.Context) error {
	participants, err := h.service.ListAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list participants")
	}
	return c.JSON(http.StatusOK, participants)
}

// FindParticipantByID godoc
// @Summary      Get a participant by ID
// @Tags         participants
// @Produce      json
// @Param        id path string true "Participant ID (UUID)"
// @Success      200 {object} models.Participant
// @Failure      400 {object} object "Invalid participant id"
// @Failure      404 {object} object "Participant not found"
// @Failure      500 {object} object "Internal error"
// @Router       /participants/{id} [get]
func (h *ParticipantHandler) FindByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid participant id")
	}

	participant, err := h.service.FindByID(id)
	if err != nil {
		if errors.Is(err, services.ErrParticipantNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find participant")
	}

	return c.JSON(http.StatusOK, participant)
}

// UpdateParticipant godoc
// @Summary      Update a participant
// @Tags         participants
// @Accept       json
// @Produce      json
// @Param        id path string true "Participant ID (UUID)"
// @Param        request body participantRequest true "Updated participant data"
// @Success      200 {object} models.Participant
// @Failure      400 {object} object "Invalid body or id, or missing name/email"
// @Failure      404 {object} object "Participant not found"
// @Failure      409 {object} object "Email already in use"
// @Failure      500 {object} object "Internal error"
// @Router       /participants/{id} [put]
func (h *ParticipantHandler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid participant id")
	}

	var req participantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	participant := &models.Participant{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.service.Update(participant); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidParticipant):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrParticipantNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrEmailTaken):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update participant")
		}
	}

	return c.JSON(http.StatusOK, participant)
}

// DeleteParticipant godoc
// @Summary      Delete a participant
// @Tags         participants
// @Param        id path string true "Participant ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} object "Invalid participant id"
// @Failure      404 {object} object "Participant not found"
// @Failure      500 {object} object "Internal error"
// @Router       /participants/{id} [delete]
func (h *ParticipantHandler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid participant id")
	}

	if err := h.service.DeleteByID(id); err != nil {
		if errors.Is(err, services.ErrParticipantNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete participant")
	}

	return c.NoContent(http.StatusNoContent)
}

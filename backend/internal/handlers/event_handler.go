package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"uuid"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/middleware"
	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/services"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

type eventRequest struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Capacity       int       `json:"capacity"`
	Duration       float64   `json:"duration"`
	StartTimestamp time.Time `json:"startTimestamp" format:"date-time"`
}

func (r eventRequest) toEvent(id uuid.UUID) *models.Event {
	return &models.Event{
		ID:             id,
		Name:           r.Name,
		Description:    r.Description,
		Location:       r.Location,
		Capacity:       r.Capacity,
		Duration:       r.Duration,
		StartTimestamp: r.StartTimestamp,
	}
}

func (h *EventHandler) parseTimestamp(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Creates an event with a name and start timestamp
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        request body eventRequest true "Event to create"
// @Success      201 {object} models.Event
// @Failure      400 {object} object "Invalid body, or missing name/startTimestamp"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events [post]
func (h *EventHandler) Create(c *echo.Context) error {
	var req eventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	event := req.toEvent(uuid.Nil())

	if err := h.service.Create(event); err != nil {
		if errors.Is(err, services.ErrInvalidEvent) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create event")
	}

	return c.JSON(http.StatusCreated, event)
}

// ListEvents godoc
// @Summary      List events with optional filters
// @Description  All filters are optional and combine; results are ordered by start time
// @Tags         events
// @Produce      json
// @Param        q query string false "Partial, case-insensitive match on name, description or location"
// @Param        after query string false "RFC3339 timestamp — only events starting after this"
// @Param        before query string false "RFC3339 timestamp — only events starting before this"
// @Param        full query boolean false "Filter by fullness: true = only full events, false = only events with space remaining"
// @Success      200 {array} models.Event
// @Failure      400 {object} object "Invalid filter value"
// @Failure      500 {object} object "Internal error"
// @Router       /api/v1/events [get]
func (h *EventHandler) List(c *echo.Context) error {
	query, err := h.parseListQuery(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	events, err := h.service.List(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list events")
	}
	return c.JSON(http.StatusOK, events)
}

func (h *EventHandler) parseListQuery(c *echo.Context) (models.EventQuery, error) {
	query := models.EventQuery{Search: c.QueryParam("q")}

	if value := c.QueryParam("after"); value != "" {
		after, err := h.parseTimestamp(value)
		if err != nil {
			return query, errors.New("invalid 'after' filter, must be RFC3339")
		}
		query.After = &after
	}

	if value := c.QueryParam("before"); value != "" {
		before, err := h.parseTimestamp(value)
		if err != nil {
			return query, errors.New("invalid 'before' filter, must be RFC3339")
		}
		query.Before = &before
	}

	if value := c.QueryParam("full"); value != "" {
		full, err := strconv.ParseBool(value)
		if err != nil {
			return query, errors.New("invalid 'full' filter, must be true or false")
		}
		query.Full = &full
	}

	return query, nil
}

// FindEventByID godoc
// @Summary      Get an event by ID with its participants
// @Tags         events
// @Produce      json
// @Param        id path string true "Event ID (UUID)"
// @Success      200 {object} models.Event
// @Failure      400 {object} object "Invalid event id"
// @Failure      404 {object} object "Event not found"
// @Failure      500 {object} object "Internal error"
// @Router       /api/v1/events/{id} [get]
func (h *EventHandler) FindByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	event, err := h.service.FindByID(id)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find event")
	}

	return c.JSON(http.StatusOK, event)
}

// JoinEvent godoc
// @Summary      Register the current user as a participant of an event
// @Description  Adds the authenticated participant to the event if there is capacity remaining and they are not already registered
// @Tags         events
// @Param        id path string true "Event ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} object "Invalid event id"
// @Failure      401 {object} object "Not authenticated"
// @Failure      404 {object} object "Event not found"
// @Failure      409 {object} object "Already registered, event is at capacity, or event has already started"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/participants [post]
func (h *EventHandler) Join(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	participant := middleware.ParticipantFrom(c)
	if participant == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}

	if err := h.service.JoinEvent(id, participant.ID); err != nil {
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrAlreadyRegistered):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, services.ErrEventFull):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, services.ErrEventStarted):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to join event")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// InviteParticipant godoc
// @Summary      Add a specific participant to an event
// @Description  Demo endpoint — adds the given participant to the event, subject to the same rules as joining (capacity, start time, duplicates)
// @Tags         events
// @Param        id path string true "Event ID (UUID)"
// @Param        participantId path string true "Participant ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} object "Invalid event or participant id"
// @Failure      401 {object} object "Not authenticated"
// @Failure      404 {object} object "Event or participant not found"
// @Failure      409 {object} object "Already registered, event is at capacity, or event has already started"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/participants/{participantId} [post]
func (h *EventHandler) Invite(c *echo.Context) error {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	participantID, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid participant id")
	}

	if middleware.ParticipantFrom(c) == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}

	if err := h.service.JoinEvent(eventID, participantID); err != nil {
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrAlreadyRegistered):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, services.ErrEventFull):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, services.ErrEventStarted):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, gorm.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "participant not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to invite participant")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// LeaveEvent godoc
// @Summary      Remove the current user's participation from an event
// @Tags         events
// @Param        id path string true "Event ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} object "Invalid event id"
// @Failure      401 {object} object "Not authenticated"
// @Failure      404 {object} object "Event not found, or not registered"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/participants [delete]
func (h *EventHandler) Leave(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	participant := middleware.ParticipantFrom(c)
	if participant == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}

	if err := h.service.LeaveEvent(id, participant.ID); err != nil {
		switch {
		case errors.Is(err, services.ErrEventNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrNotRegistered):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to leave event")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// FindEventByName godoc
// @Summary      Find an event by exact name
// @Tags         events
// @Produce      json
// @Param        name path string true "Event name"
// @Success      200 {object} models.Event
// @Failure      404 {object} object "Event not found"
// @Failure      500 {object} object "Internal error"
// @Router       /api/v1/events/name/{name} [get]
func (h *EventHandler) FindByName(c *echo.Context) error {
	name := c.Param("name")

	event, err := h.service.FindByName(name)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find event")
	}

	return c.JSON(http.StatusOK, event)
}

// SearchEvents godoc
// @Summary      Search events by partial name, description, or location
// @Tags         events
// @Produce      json
// UpdateEvent godoc
// @Summary      Update an event
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id path string true "Event ID (UUID)"
// @Param        request body eventRequest true "Updated event data"
// @Success      200 {object} models.Event
// @Failure      400 {object} object "Invalid body or id, or missing name/startTimestamp"
// @Failure      404 {object} object "Event not found"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [put]
func (h *EventHandler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	var req eventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	event := req.toEvent(id)

	if err := h.service.Update(event); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidEvent):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrEventNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update event")
		}
	}

	return c.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Delete an event
// @Tags         events
// @Param        id path string true "Event ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} object "Invalid event id"
// @Failure      404 {object} object "Event not found"
// @Failure      500 {object} object "Internal error"
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [delete]
func (h *EventHandler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	}

	if err := h.service.DeleteByID(id); err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete event")
	}

	return c.NoContent(http.StatusNoContent)
}

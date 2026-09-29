package handlers

import (
	"errors"
	"net/http"
	"time"
	"uuid"

	"github.com/labstack/echo/v5"

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
// @Router       /events [post]
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
// @Summary      List all events
// @Tags         events
// @Produce      json
// @Success      200 {array} models.Event
// @Failure      500 {object} object "Internal error"
// @Router       /events [get]
func (h *EventHandler) List(c *echo.Context) error {
	events, err := h.service.ListAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list events")
	}
	return c.JSON(http.StatusOK, events)
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
// @Router       /events/{id} [get]
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

// FindEventByName godoc
// @Summary      Find an event by exact name
// @Tags         events
// @Produce      json
// @Param        name path string true "Event name"
// @Success      200 {object} models.Event
// @Failure      404 {object} object "Event not found"
// @Failure      500 {object} object "Internal error"
// @Router       /events/name/{name} [get]
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
// @Param        q query string true "Search term"
// @Success      200 {array} models.Event
// @Failure      400 {object} object "Missing search term"
// @Failure      500 {object} object "Internal error"
// @Router       /events/search [get]
func (h *EventHandler) Search(c *echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing search term")
	}

	events, err := h.service.FindByPartialNameDescriptionLocation(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to search events")
	}
	return c.JSON(http.StatusOK, events)
}

// FindEventsBefore godoc
// @Summary      Find events starting before a timestamp
// @Tags         events
// @Produce      json
// @Param        timestamp query string true "RFC3339 timestamp"
// @Success      200 {array} models.Event
// @Failure      400 {object} object "Invalid or missing timestamp"
// @Failure      500 {object} object "Internal error"
// @Router       /events/before [get]
func (h *EventHandler) FindAllBefore(c *echo.Context) error {
	timestamp, err := h.parseTimestamp(c.QueryParam("timestamp"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or missing timestamp, must be RFC3339")
	}

	events, err := h.service.FindAllBefore(timestamp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find events")
	}
	return c.JSON(http.StatusOK, events)
}

// FindEventsAfter godoc
// @Summary      Find events starting after a timestamp
// @Tags         events
// @Produce      json
// @Param        timestamp query string true "RFC3339 timestamp"
// @Success      200 {array} models.Event
// @Failure      400 {object} object "Invalid or missing timestamp"
// @Failure      500 {object} object "Internal error"
// @Router       /events/after [get]
func (h *EventHandler) FindAllAfter(c *echo.Context) error {
	timestamp, err := h.parseTimestamp(c.QueryParam("timestamp"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or missing timestamp, must be RFC3339")
	}

	events, err := h.service.FindAllAfter(timestamp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find events")
	}
	return c.JSON(http.StatusOK, events)
}

// FindEventsBetween godoc
// @Summary      Find events starting between two timestamps
// @Tags         events
// @Produce      json
// @Param        start query string true "RFC3339 start timestamp"
// @Param        end query string true "RFC3339 end timestamp"
// @Success      200 {array} models.Event
// @Failure      400 {object} object "Invalid or missing timestamps"
// @Failure      500 {object} object "Internal error"
// @Router       /events/between [get]
func (h *EventHandler) FindAllBetween(c *echo.Context) error {
	start, err := h.parseTimestamp(c.QueryParam("start"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or missing start timestamp, must be RFC3339")
	}

	end, err := h.parseTimestamp(c.QueryParam("end"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or missing end timestamp, must be RFC3339")
	}

	events, err := h.service.FindAllBetween(start, end)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find events")
	}
	return c.JSON(http.StatusOK, events)
}

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
// @Router       /events/{id} [put]
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
// @Router       /events/{id} [delete]
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

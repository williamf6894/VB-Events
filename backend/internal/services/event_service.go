package services

import (
	"errors"
	"log/slog"
	"time"
	"uuid"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidEvent  = errors.New("invalid event data")
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) Create(event *models.Event) error {
	if event.Name == "" || event.StartTimestamp.IsZero() {
		return ErrInvalidEvent
	}
	return s.repo.Create(event)
}

func (s *EventService) Update(event *models.Event) error {
	if event.Name == "" || event.StartTimestamp.IsZero() {
		return ErrInvalidEvent
	}

	_, err := s.repo.FindByID(event.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEventNotFound
		}
		return err
	}

	return s.repo.Update(event)
}

func (s *EventService) DeleteByID(id uuid.UUID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEventNotFound
		}
		return err
	}
	return nil
}

func (s *EventService) ListAll() ([]models.Event, error) {
	events, err := s.repo.ListAll()
	if err != nil {
		slog.Error("failed to list events", "error", err)
		return nil, err
	}
	return events, nil
}

func (s *EventService) FindByID(id uuid.UUID) (*models.Event, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (s *EventService) FindByName(name string) (*models.Event, error) {
	event, err := s.repo.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (s *EventService) FindByPartialNameDescriptionLocation(query string) ([]models.Event, error) {
	events, err := s.repo.FindByPartialNameDescriptionLocation(query)
	if err != nil {
		slog.Error("failed to search events", "error", err, "query", query)
		return nil, err
	}
	return events, nil
}

func (s *EventService) FindAllBefore(timestamp time.Time) ([]models.Event, error) {
	events, err := s.repo.FindAllBefore(timestamp)
	if err != nil {
		slog.Error("failed to find events before timestamp", "error", err, "timestamp", timestamp)
		return nil, err
	}
	return events, nil
}

func (s *EventService) FindAllAfter(timestamp time.Time) ([]models.Event, error) {
	events, err := s.repo.FindAllAfter(timestamp)
	if err != nil {
		slog.Error("failed to find events after timestamp", "error", err, "timestamp", timestamp)
		return nil, err
	}
	return events, nil
}

func (s *EventService) FindAllBetween(start, end time.Time) ([]models.Event, error) {
	events, err := s.repo.FindAllBetween(start, end)
	if err != nil {
		slog.Error("failed to find events between timestamps", "error", err, "start", start, "end", end)
		return nil, err
	}
	return events, nil
}

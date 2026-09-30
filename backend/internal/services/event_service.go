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
	ErrEventNotFound     = errors.New("event not found")
	ErrInvalidEvent      = errors.New("invalid event data")
	ErrAlreadyRegistered = errors.New("already registered for this event")
	ErrEventFull         = errors.New("event is at capacity")
	ErrEventStarted      = errors.New("event has already started")
	ErrNotRegistered     = errors.New("not registered for this event")
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

func (s *EventService) List(query models.EventQuery) ([]models.Event, error) {
	events, err := s.repo.List(query)
	if err != nil {
		slog.Error("failed to list events", "error", err, "search", query.Search)
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

func (s *EventService) JoinEvent(eventID, participantID uuid.UUID) error {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEventNotFound
		}
		return err
	}

	if !event.StartTimestamp.After(time.Now()) {
		return ErrEventStarted
	}

	for _, participant := range event.Participants {
		if participant.ID == participantID {
			return ErrAlreadyRegistered
		}
	}

	if event.Capacity > 0 && len(event.Participants) >= event.Capacity {
		return ErrEventFull
	}

	return s.repo.AddParticipant(eventID, participantID)
}

func (s *EventService) LeaveEvent(eventID, participantID uuid.UUID) error {
	_, err := s.repo.FindByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEventNotFound
		}
		return err
	}

	if err := s.repo.RemoveParticipant(eventID, participantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotRegistered
		}
		return err
	}

	return nil
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

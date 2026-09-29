package repository

import (
	"log/slog"
	"time"
	"uuid"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	if err := r.db.Create(event).Error; err != nil {
		slog.Error("failed to create event", "error", err)
		return err
	}
	return nil
}

func (r *EventRepository) Update(event *models.Event) error {
	if err := r.db.Save(event).Error; err != nil {
		slog.Error("failed to update event", "error", err, "id", event.ID)
		return err
	}
	return nil
}

func (r *EventRepository) DeleteByID(id uuid.UUID) error {
	if err := r.db.Exec("DELETE FROM event_participants WHERE event_id = ?", id).Error; err != nil {
		slog.Error("failed to remove event registrations", "error", err, "id", id)
		return err
	}

	result := r.db.Delete(&models.Event{}, "id = ?", id)
	if result.Error != nil {
		slog.Error("failed to delete event", "error", result.Error, "id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *EventRepository) ListAll() ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Order("start_timestamp").Find(&events).Error; err != nil {
		slog.Error("failed to list events", "error", err)
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) FindByID(id uuid.UUID) (*models.Event, error) {
	var event models.Event
	if err := r.db.Preload("Participants").First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) FindByName(name string) (*models.Event, error) {
	var event models.Event
	if err := r.db.Where("name = ?", name).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) AddParticipant(eventID, participantID uuid.UUID) error {
	var event models.Event
	if err := r.db.First(&event, "id = ?", eventID).Error; err != nil {
		return err
	}

	var participant models.Participant
	if err := r.db.First(&participant, "id = ?", participantID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&event).Association("Participants").Append(&participant); err != nil {
		slog.Error("failed to add participant to event", "error", err, "eventID", eventID, "participantID", participantID)
		return err
	}
	return nil
}

func (r *EventRepository) RemoveParticipant(eventID, participantID uuid.UUID) error {
	result := r.db.Exec("DELETE FROM event_participants WHERE event_id = ? AND participant_id = ?", eventID, participantID)
	if result.Error != nil {
		slog.Error("failed to remove participant from event", "error", result.Error, "eventID", eventID, "participantID", participantID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *EventRepository) FindByPartialNameDescriptionLocation(query string) ([]models.Event, error) {
	var events []models.Event
	pattern := "%" + query + "%"
	if err := r.db.Where("name ILIKE ? OR description ILIKE ? OR location ILIKE ?", pattern, pattern, pattern).
		Order("start_timestamp").
		Find(&events).Error; err != nil {
		slog.Error("failed to search events", "error", err, "query", query)
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) FindAllBefore(timestamp time.Time) ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Where("start_timestamp < ?", timestamp).
		Order("start_timestamp").
		Find(&events).Error; err != nil {
		slog.Error("failed to find events before timestamp", "error", err, "timestamp", timestamp)
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) FindAllAfter(timestamp time.Time) ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Where("start_timestamp > ?", timestamp).
		Order("start_timestamp").
		Find(&events).Error; err != nil {
		slog.Error("failed to find events after timestamp", "error", err, "timestamp", timestamp)
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) FindAllBetween(start, end time.Time) ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Where("start_timestamp BETWEEN ? AND ?", start, end).
		Order("start_timestamp").
		Find(&events).Error; err != nil {
		slog.Error("failed to find events between timestamps", "error", err, "start", start, "end", end)
		return nil, err
	}
	return events, nil
}

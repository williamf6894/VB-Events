package repository

import (
	"log/slog"

	"gorm.io/gorm"
	"uuid"

	"github.com/williamf6894/VB-Events/internal/models"
)

type ParticipantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) Create(participant *models.Participant) error {
	if err := r.db.Create(participant).Error; err != nil {
		slog.Error("failed to create participant", "error", err)
		return err
	}
	return nil
}

func (r *ParticipantRepository) FindByID(id uuid.UUID) (*models.Participant, error) {
	var participant models.Participant
	if err := r.db.First(&participant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *ParticipantRepository) FindByEmail(email string) (*models.Participant, error) {
	var participant models.Participant
	if err := r.db.Where("email = ?", email).First(&participant).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *ParticipantRepository) Update(participant *models.Participant) error {
	if err := r.db.Save(participant).Error; err != nil {
		slog.Error("failed to update participant", "error", err, "id", participant.ID)
		return err
	}
	return nil
}

func (r *ParticipantRepository) DeleteByID(id uuid.UUID) error {
	result := r.db.Delete(&models.Participant{}, "id = ?", id)
	if result.Error != nil {
		slog.Error("failed to delete participant", "error", result.Error, "id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ParticipantRepository) ListAll() ([]models.Participant, error) {
	var participants []models.Participant
	if err := r.db.Find(&participants).Error; err != nil {
		slog.Error("failed to list participants", "error", err)
		return nil, err
	}
	return participants, nil
}

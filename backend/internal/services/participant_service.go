package services

import (
	"errors"
	"log/slog"

	"uuid"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
)

const minPasswordLength = 8

var (
	ErrParticipantNotFound = errors.New("participant not found")
	ErrEmailTaken          = errors.New("participant email already in use")
	ErrInvalidParticipant  = errors.New("invalid participant data")
)

type ParticipantService struct {
	repo *repository.ParticipantRepository
}

func NewParticipantService(repo *repository.ParticipantRepository) *ParticipantService {
	return &ParticipantService{repo: repo}
}

func (s *ParticipantService) FindByID(id uuid.UUID) (*models.Participant, error) {
	participant, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParticipantNotFound
		}
		return nil, err
	}
	return participant, nil
}

func (s *ParticipantService) FindByEmail(email string) (*models.Participant, error) {
	participant, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParticipantNotFound
		}
		return nil, err
	}
	return participant, nil
}

func (s *ParticipantService) Update(participant *models.Participant) error {
	if !validParticipant(participant) {
		return ErrInvalidParticipant
	}

	current, err := s.repo.FindByID(participant.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrParticipantNotFound
		}
		return err
	}

	if current.Email != participant.Email {
		existing, err := s.repo.FindByEmail(participant.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existing != nil {
			return ErrEmailTaken
		}
	}

	hash, err := hashPassword(participant.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return err
	}
	participant.Password = hash

	return s.repo.Update(participant)
}

func (s *ParticipantService) DeleteByID(id uuid.UUID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrParticipantNotFound
		}
		return err
	}
	return nil
}

func (s *ParticipantService) ListAll() ([]models.Participant, error) {
	participants, err := s.repo.ListAll()
	if err != nil {
		slog.Error("failed to list participants", "error", err)
		return nil, err
	}
	return participants, nil
}

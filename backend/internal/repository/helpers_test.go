package repository

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/williamf6894/VB-Events/internal/models"
)

func cleanTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("TRUNCATE events, participants, event_participants CASCADE").Error; err != nil {
		t.Fatalf("failed to truncate tables: %s", err)
	}
}

func newParticipant(name, email string) *models.Participant {
	return &models.Participant{
		Name:  name,
		Email: email,
	}
}

func newEvent(name string, startTimestamp time.Time) *models.Event {
	return &models.Event{
		Name:           name,
		Description:    name + " description",
		Location:       name + " location",
		Capacity:       10,
		Duration:       1.5,
		StartTimestamp: startTimestamp,
	}
}

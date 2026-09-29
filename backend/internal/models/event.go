package models

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Event struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id" swaggertype:"string"`
	Name           string        `gorm:"type:text;not null" json:"name"`
	Description    string        `gorm:"type:text" json:"description"`
	Location       string        `gorm:"type:text" json:"location"`
	Capacity       int           `gorm:"type:int" json:"capacity"`
	Duration       float64       `gorm:"type:double precision" json:"duration"`
	StartTimestamp time.Time     `gorm:"type:timestamptz;not null" json:"startTimestamp"`
	Participants   []Participant `gorm:"many2many:event_participants;" json:"participants,omitempty"`
}

// Specifically want to use UUIDV7 for timestamp sorting.
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil() {
		e.ID = uuid.NewV7()
	}
	return nil
}

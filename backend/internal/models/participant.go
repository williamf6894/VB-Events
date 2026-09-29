package models

import (
	"uuid"

	"gorm.io/gorm"
)

type Participant struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id" swaggertype:"string"`
	Name     string    `gorm:"type:text;not null" json:"name"`
	Email    string    `gorm:"type:text;not null;uniqueIndex" json:"email"`
	Password string    `gorm:"type:text" json:"-"`
}

func (p *Participant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil() {
		p.ID = uuid.NewV7()
	}
	return nil
}

package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Model struct {
	Id        string `gorm:"primary_key;size:32"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (m *Model) BeforeCreate(_ *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	m.Id = uuid.String()

	return nil
}
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseEntity is a common structure for database entities with basic fields like ID, timestamps, and soft delete support.
type BaseEntity struct {
	Id        *uuid.UUID `gorm:"primary_key" json:"id" sql:"index"`
	CreatedAt time.Time  `json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (b *BaseEntity) BeforeCreate(_ *gorm.DB) error {
	newUUID := uuid.New()
	b.Id = &newUUID
	return nil
}

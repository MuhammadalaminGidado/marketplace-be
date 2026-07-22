package models

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email string `gorm:"uniqueIndex;not null"`

	PasswordDigest string `gorm:"not null"`

	Status string `gorm:"type:entity_status;default:'active';not null"`

	Sessions []Session `gorm:"foreignKey:EntityID"`

	LastLoginAt *time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

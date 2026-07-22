package models

import "time"

type Session struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	EntityID string `gorm:"type:uuid;not null;index"`

	TokenDigest string `gorm:"column:token_digest;uniqueIndex;not null"`

	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time
}

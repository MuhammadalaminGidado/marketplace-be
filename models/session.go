package models

import "time"

type Session struct {
	ID        string    `json:"id,omitempty"    gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"userId"          gorm:"column:user_id;type:uuid;not null"`
	Token     string    `json:"token"    gorm:"column:token;uniqueIndex;not null"`
	CreatedAt time.Time `json:"createdAt"       gorm:"autoCreateTime"`
	ExpiresAt time.Time `json:"expiresAt"       gorm:"not null"`
}

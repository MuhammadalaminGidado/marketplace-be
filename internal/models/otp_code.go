package models

import "time"

const (
	OtpPurposeVerifyEmail   = "verify_email"
	OtpPurposeResetPassword = "reset_password"
)

type OTPCode struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	EntityID string `gorm:"type:uuid;not null;index"`

	Purpose string `gorm:"not null"`

	CodeDigest string `gorm:"not null"`

	ConsumedAt *time.Time

	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time
}

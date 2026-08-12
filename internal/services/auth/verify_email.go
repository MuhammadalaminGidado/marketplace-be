package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example/api/internal/auth"
	"example/api/internal/models"

	"gorm.io/gorm"
)

func (s *Service) VerifyEmail(
	ctx context.Context,
	email, code string,
) error {

	entity, err := s.entities.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidOrExpiredCode
		}
		return fmt.Errorf("find entity: %w", err)
	}

	otp, err := s.otps.FindUnconsumed(
		ctx,
		entity.ID,
		models.OtpPurposeVerifyEmail,
		auth.HashOTP(code),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidOrExpiredCode
		}
		return fmt.Errorf("find otp: %w", err)
	}

	if time.Now().After(otp.ExpiresAt) {
		return ErrInvalidOrExpiredCode
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		otpRepo := s.otps.WithDB(tx)
		entityRepo := s.entities.WithDB(tx)

		if err := otpRepo.Consume(ctx, otp.ID); err != nil {
			return ErrInvalidOrExpiredCode
		}

		if err := entityRepo.MarkEmailVerified(ctx, entity.ID); err != nil {
			return fmt.Errorf("mark email verified: %w", err)
		}

		return nil
	})
}

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

func (s *Service) ResetPassword(
	ctx context.Context,
	email, code, newPassword string,
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
		models.OtpPurposeResetPassword,
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

	digest, err := auth.PasswordHash(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		otpRepo := s.otps.WithDB(tx)
		entityRepo := s.entities.WithDB(tx)
		sessionRepo := s.sessions.WithDB(tx)

		if err := otpRepo.Consume(ctx, otp.ID); err != nil {
			return ErrInvalidOrExpiredCode
		}

		if err := entityRepo.UpdatePasswordDigest(ctx, entity.ID, digest); err != nil {
			return fmt.Errorf("update password digest: %w", err)
		}

		if err := sessionRepo.DeleteByEntityID(ctx, entity.ID); err != nil {
			return fmt.Errorf("delete sessions: %w", err)
		}

		return nil
	})
}

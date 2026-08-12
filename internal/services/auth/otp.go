package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"example/api/internal/auth"
	"example/api/internal/mailer"
	"example/api/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) RequestOTP(
	ctx context.Context,
	email, purpose string,
) error {

	if err := s.throttleOTP(email); err != nil {
		return err
	}

	entity, err := s.entities.FindByEmail(ctx, email)
	if err != nil {
		// Leak-nothing: emit the same success response whether or not the
		// account exists.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("find entity: %w", err)
	}

	s.issueAndSendOTP(ctx, entity, purpose)

	return nil
}

func (s *Service) issueAndSendOTP(
	ctx context.Context,
	entity *models.Entity,
	purpose string,
) {

	code, err := auth.GenerateOTP()
	if err != nil {
		s.logger.Error("generate otp", zap.Error(err))
		return
	}

	otp := &models.OTPCode{
		EntityID:   entity.ID,
		Purpose:    purpose,
		CodeDigest: auth.HashOTP(code),
		ExpiresAt:  time.Now().Add(otpTTL),
	}

	if err := s.otps.IssueActive(ctx, otp); err != nil {
		s.logger.Error("issue otp", zap.Error(err))
		return
	}

	appName := s.mailer.AppName()
	subject := appName + " verification code"
	body := fmt.Sprintf(
		"Your %s code is %s. It expires in 5 minutes. If you didn't request this, you can ignore this email.",
		appName,
		code,
	)

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), mailer.SendTimeout()+15*time.Second)
		defer cancel()
		if err := s.mailer.Send(bgCtx, entity.Email, subject, body); err != nil {
			s.logger.Error("send otp email", zap.String("to", entity.Email), zap.Error(err))
		}
	}()
}

func (s *Service) throttleOTP(email string) error {
	if s.redis == nil {
		return nil
	}

	key := "ratelimit:otp:" + strings.ToLower(email)

	ctx := context.Background()

	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		// Redis down — fail open so an outage doesn't block OTP requests.
		s.logger.Warn("otp rate limiter unavailable", zap.Error(err))
		return nil
	}

	if count == 1 {
		s.redis.Expire(ctx, key, time.Hour)
	}

	if count > otpMaxPerHour {
		return ErrOtpThrottled
	}

	return nil
}

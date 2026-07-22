package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example/api/internal/auth"
	"example/api/internal/dto"
	"example/api/internal/models"

	"gorm.io/gorm"
)

func (s *Service) Signup(
	ctx context.Context,
	req dto.SignupPayload,
) (*models.Entity, string, string, error) {

	if req.Password != req.PasswordConfirmation {
		return nil, "", "", ErrPasswordMismatch
	}

	passwordDigest, err := auth.PasswordHash(req.Password)
	if err != nil {
		return nil, "", "", fmt.Errorf("hash password: %w", err)
	}

	entity := &models.Entity{
		Email:          req.Email,
		PasswordDigest: passwordDigest,
	}

	var token string
	var csrfToken string
	err = s.db.Transaction(func(tx *gorm.DB) error {
		entityRepo := s.entities.WithDB(tx)
		sessionRepo := s.sessions.WithDB(tx)

		if err := entityRepo.Create(ctx, entity); err != nil {

			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrEmailAlreadyExists
			}

			return err
		}

		var tokenDigest string
		token, tokenDigest, err = auth.GenerateSessionToken()
		if err != nil {
			return err
		}

		csrfToken, err = auth.GenerateCSRFToken()
		if err != nil {
			return fmt.Errorf("generate csrf token: %w", err)
		}

		session := &models.Session{
			EntityID:    entity.ID,
			TokenDigest: tokenDigest,
			ExpiresAt:   time.Now().Add(s.sessionDuration),
		}

		return sessionRepo.Create(ctx, session)
	})

	if err != nil {
		return nil, "", "", err
	}

	return entity, token, csrfToken, nil
}

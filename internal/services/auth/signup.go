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
	passwordDigest, err := auth.PasswordHash(req.Password)
	if err != nil {
		return nil, "", "", fmt.Errorf("hash password: %w", err)
	}

	entity := &models.Entity{
		Email:          req.Email,
		PasswordDigest: passwordDigest,
	}

	token, tokenDigest, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("generate session token: %w", err)
	}
	csrfToken, err := auth.GenerateCSRFToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("generate csrf token: %w", err)
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		entityRepo := s.entities.WithDB(tx)
		sessionRepo := s.sessions.WithDB(tx)

		if err := entityRepo.Create(ctx, entity); err != nil {

			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrEmailAlreadyExists
			}

			return err
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

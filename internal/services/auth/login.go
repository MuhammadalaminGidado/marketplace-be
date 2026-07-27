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

func (s *Service) Login(
	ctx context.Context,
	req dto.LoginPayload,
) (*models.Entity, string, string, error) {

	entity, err := s.entities.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			auth.CompareWithDummyHash(req.Password)
			return nil, "", "", ErrInvalidCredentials
		}

		return nil, "", "", fmt.Errorf("find entity: %w", err)
	}

	if !auth.CheckPasswordHash(
		entity.PasswordDigest,
		req.Password,
	) {
		return nil, "", "", ErrInvalidCredentials
	}

	if entity.Status != models.StatusActive {
		return nil, "", "", ErrInvalidCredentials
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

		if err := entityRepo.UpdateLastLogin(ctx, entity.ID); err != nil {
			return fmt.Errorf("update last login: %w", err)
		}

		session := &models.Session{
			EntityID:    entity.ID,
			TokenDigest: tokenDigest,
			ExpiresAt:   time.Now().Add(s.sessionDuration),
		}

		if err := sessionRepo.Create(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, "", "", err
	}

	return entity, token, csrfToken, nil
}

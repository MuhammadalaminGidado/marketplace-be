package auth

import (
	"context"

	"example/api/internal/models"
)

func (s *Service) CurrentEntity(
	ctx context.Context,
	token string,
) (*models.Entity, error) {

	entityID, err := s.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.entities.FindByID(ctx, entityID)
}

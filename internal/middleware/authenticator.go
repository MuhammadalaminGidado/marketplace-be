package middleware

import (
	"context"

	"example/api/internal/models"
)

type Authenticator interface {
	CurrentEntity(ctx context.Context, token string) (*models.Entity, error)
}

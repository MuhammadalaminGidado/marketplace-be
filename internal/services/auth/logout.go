package auth

import (
	"context"
	"example/api/internal/auth"
)

func (s *Service) Logout(
	ctx context.Context,
	token string,
) error {
	// logout.go
	return s.sessions.DeleteByTokenDigest(ctx, auth.HashSessionToken(token))
}

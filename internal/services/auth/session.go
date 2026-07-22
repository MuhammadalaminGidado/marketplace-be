package auth

import (
	"context"
	"example/api/internal/auth"
	"time"
)

func (s *Service) ValidateSession(
	ctx context.Context,
	token string,
) (string, error) {

	session, err := s.sessions.FindByTokenDigest(ctx, auth.HashSessionToken(token))
	if err != nil {
		return "", err
	}

	if time.Now().After(session.ExpiresAt) {

		_ = s.sessions.DeleteByTokenDigest(ctx, auth.HashSessionToken(token))

		return "", ErrSessionExpired
	}

	return session.EntityID, nil
}

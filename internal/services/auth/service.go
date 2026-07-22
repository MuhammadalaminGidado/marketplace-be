package auth

import (
	"time"

	"example/api/internal/repositories"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db              *gorm.DB
	entities        repositories.EntityRepository
	sessions        repositories.SessionRepository
	logger          *zap.Logger
	sessionDuration time.Duration
	maxSessions     int
}

func NewService(
	db *gorm.DB,
	entities repositories.EntityRepository,
	sessions repositories.SessionRepository,
	logger *zap.Logger,
	sessionDuration time.Duration,
	maxSessions int,
) *Service {
	return &Service{
		db:              db,
		entities:        entities,
		sessions:        sessions,
		logger:          logger,
		sessionDuration: sessionDuration,
		maxSessions:     maxSessions,
	}
}

package auth

import (
	"time"

	"example/api/internal/mailer"
	"example/api/internal/repositories"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	otpTTL        = 5 * time.Minute
	otpMaxPerHour = 3
)

type Service struct {
	db              *gorm.DB
	entities        repositories.EntityRepository
	sessions        repositories.SessionRepository
	otps            repositories.OTPRepository
	mailer          mailer.Mailer
	redis           *redis.Client
	logger          *zap.Logger
	sessionDuration time.Duration
	maxSessions     int
}

func NewService(
	db *gorm.DB,
	entities repositories.EntityRepository,
	sessions repositories.SessionRepository,
	otps repositories.OTPRepository,
	mailer mailer.Mailer,
	redis *redis.Client,
	logger *zap.Logger,
	sessionDuration time.Duration,
	maxSessions int,
) *Service {
	return &Service{
		db:              db,
		entities:        entities,
		sessions:        sessions,
		otps:            otps,
		mailer:          mailer,
		redis:           redis,
		logger:          logger,
		sessionDuration: sessionDuration,
		maxSessions:     maxSessions,
	}
}

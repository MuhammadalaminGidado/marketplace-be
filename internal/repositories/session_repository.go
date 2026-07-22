package repositories

import (
	"context"

	"example/api/internal/models"

	"gorm.io/gorm"
)

type SessionRepository interface {
	WithDB(db *gorm.DB) SessionRepository

	Create(ctx context.Context, session *models.Session) error

	FindByTokenDigest(ctx context.Context, tokenDigest string) (*models.Session, error)

	Delete(ctx context.Context, id string) error
	DeleteByTokenDigest(ctx context.Context, tokenDigest string) error
	DeleteByEntityID(ctx context.Context, entityID string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func (r *sessionRepository) WithDB(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(
	ctx context.Context,
	session *models.Session,
) error {
	return r.db.
		WithContext(ctx).
		Create(session).
		Error
}

func (r *sessionRepository) FindByTokenDigest(
	ctx context.Context,
	tokenDigest string,
) (*models.Session, error) {

	var session models.Session

	err := r.db.
		WithContext(ctx).
		Where("token_digest = ?", tokenDigest).
		First(&session).
		Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) Delete(
	ctx context.Context,
	id string,
) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.Session{}, "id = ?", id).
		Error
}

func (r *sessionRepository) DeleteByTokenDigest(
	ctx context.Context,
	tokenDigest string,
) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.Session{}, "token_digest = ?", tokenDigest).
		Error
}

func (r *sessionRepository) DeleteByEntityID(
	ctx context.Context,
	entityID string,
) error {
	return r.db.
		WithContext(ctx).
		Delete(&models.Session{}, "entity_id = ?", entityID).
		Error
}

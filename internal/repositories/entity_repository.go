package repositories

import (
	"context"

	"example/api/internal/models"

	"gorm.io/gorm"
)

type EntityRepository interface {
	WithDB(db *gorm.DB) EntityRepository

	Create(ctx context.Context, entity *models.Entity) error

	FindByEmail(ctx context.Context, email string) (*models.Entity, error)
	FindByID(ctx context.Context, id string) (*models.Entity, error)

	Update(ctx context.Context, entity *models.Entity) error
	UpdateLastLogin(ctx context.Context, entityID string) error
}

type entityRepository struct {
	db *gorm.DB
}

func (r *entityRepository) WithDB(db *gorm.DB) EntityRepository {
	return &entityRepository{
		db: db,
	}
}

func NewEntityRepository(db *gorm.DB) EntityRepository {
	return &entityRepository{
		db: db,
	}
}

func (r *entityRepository) Create(
	ctx context.Context,
	entity *models.Entity,
) error {
	return r.db.
		WithContext(ctx).
		Create(entity).
		Error
}

func (r *entityRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.Entity, error) {

	var entity models.Entity

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *entityRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Entity, error) {

	var entity models.Entity

	err := r.db.
		WithContext(ctx).Where("id = ?", id).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}
func (r *entityRepository) Update(
	ctx context.Context,
	entity *models.Entity,
) error {
	return r.db.
		WithContext(ctx).
		Save(entity).
		Error
}

func (r *entityRepository) UpdateLastLogin(
	ctx context.Context,
	entityID string,
) error {
	return r.db.
		WithContext(ctx).
		Model(&models.Entity{}).
		Where("id = ?", entityID).
		Update("last_login_at", gorm.Expr("NOW()")).
		Error
}

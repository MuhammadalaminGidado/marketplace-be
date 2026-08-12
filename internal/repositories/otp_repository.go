package repositories

import (
	"context"

	"example/api/internal/models"

	"gorm.io/gorm"
)

type OTPRepository interface {
	WithDB(db *gorm.DB) OTPRepository

	Create(ctx context.Context, otp *models.OTPCode) error
	IssueActive(ctx context.Context, otp *models.OTPCode) error
	FindUnconsumed(ctx context.Context, entityID, purpose, codeDigest string) (*models.OTPCode, error)
	Consume(ctx context.Context, id string) error
	DeleteUnconsumed(ctx context.Context, entityID, purpose string) error
}

type otpRepository struct {
	db *gorm.DB
}

func (r *otpRepository) WithDB(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) DeleteUnconsumed(
	ctx context.Context,
	entityID, purpose string,
) error {
	return r.db.
		WithContext(ctx).
		Where("entity_id = ? AND purpose = ? AND consumed_at IS NULL", entityID, purpose).
		Delete(&models.OTPCode{}).
		Error
}

func (r *otpRepository) IssueActive(
	ctx context.Context,
	otp *models.OTPCode,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		repo := r.WithDB(tx)

		if err := repo.DeleteUnconsumed(ctx, otp.EntityID, otp.Purpose); err != nil {
			return err
		}

		return repo.Create(ctx, otp)
	})
}

func (r *otpRepository) Create(
	ctx context.Context,
	otp *models.OTPCode,
) error {
	return r.db.
		WithContext(ctx).
		Create(otp).
		Error
}

func (r *otpRepository) FindUnconsumed(
	ctx context.Context,
	entityID, purpose, codeDigest string,
) (*models.OTPCode, error) {

	var otp models.OTPCode

	err := r.db.
		WithContext(ctx).
		Where(
			"entity_id = ? AND purpose = ? AND code_digest = ? AND consumed_at IS NULL",
			entityID,
			purpose,
			codeDigest,
		).
		First(&otp).
		Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *otpRepository) Consume(
	ctx context.Context,
	id string,
) error {
	res := r.db.
		WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("id = ? AND consumed_at IS NULL", id).
		Update("consumed_at", gorm.Expr("NOW()"))

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

package impl

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"medods_test_task/internal/model"
	"medods_test_task/internal/repository/intf"
)

type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) intf.RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{db: db}
}

func (r *RefreshTokenRepositoryImpl) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) GetByID(refreshTokenID uuid.UUID) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.
		Where("id = ?", refreshTokenID).
		Where("deactivated_at IS NULL").
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepositoryImpl) MarkAsDeactivated(token *model.RefreshToken) error {
	now := time.Now()
	token.DeactivatedAt = &now
	return r.db.Save(token).Error
}

func (r *RefreshTokenRepositoryImpl) MarkAllAsDeactivatedByUserID(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.RefreshToken{}).
		Where("user_id = ? AND deactivated_at IS NULL", userID).
		Update("deactivated_at", now).
		Error
}

func (r *RefreshTokenRepositoryImpl) IsTokenActive(refreshTokenID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.RefreshToken{}).
		Where("id = ? AND deactivated_at IS NULL", refreshTokenID).
		Count(&count).Error
	return count > 0, err
}

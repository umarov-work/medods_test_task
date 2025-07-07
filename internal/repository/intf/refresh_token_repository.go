package intf

import (
	"github.com/google/uuid"

	"medods_test_task/internal/model"
)

type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	GetByID(refreshTokenID uuid.UUID) (*model.RefreshToken, error)
	MarkAsDeactivated(token *model.RefreshToken) error
	MarkAllAsDeactivatedByUserID(userID uuid.UUID) error
	IsTokenActive(refreshTokenID uuid.UUID) (bool, error)
}

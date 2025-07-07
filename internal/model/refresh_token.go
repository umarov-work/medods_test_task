package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash     string     `gorm:"not null"`
	UserAgent     string     `gorm:"not null"`
	IP            string     `gorm:"not null"`
	DeactivatedAt *time.Time `gorm:"index"`
	CreatedAt     time.Time
}

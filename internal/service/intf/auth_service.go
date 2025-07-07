package intf

import "github.com/google/uuid"

type AuthService interface {
	CreateTokens(userID uuid.UUID, userAgent, ip string) (string, string, error)
	UpdateTokens(refreshTokenID uuid.UUID, rawRefreshToken, userAgent, ip string) (string, string, error)
	DeauthorizeUser(userID uuid.UUID) error
	GetUserIDByRefreshTokenID(refreshTokenID uuid.UUID) (uuid.UUID, error)
	IsTokenValid(refreshTokenID uuid.UUID) (bool, error)
}

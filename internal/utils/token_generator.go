package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"medods_test_task/internal/config"
)

func GenerateAccessToken(refreshTokenID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"refresh_token_id": refreshTokenID.String(),
		"exp":              time.Now().Add(config.Load().AccessTokenTTL).Unix(),
		"iat":              time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(config.Load().JWTSecret)
}

func GenerateRefreshToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(raw)
	return token, nil
}

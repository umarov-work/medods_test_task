package impl

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/google/uuid"

	"medods_test_task/internal/model"
	repoIntf "medods_test_task/internal/repository/intf"
	serviceIntf "medods_test_task/internal/service/intf"
	"medods_test_task/internal/utils"
)

type AuthServiceImpl struct {
	refreshTokenRepository repoIntf.RefreshTokenRepository
}

func NewAuthService(refreshTokenRepository repoIntf.RefreshTokenRepository) serviceIntf.AuthService {
	return &AuthServiceImpl{refreshTokenRepository: refreshTokenRepository}
}

func (s *AuthServiceImpl) CreateTokens(userID uuid.UUID, userAgent, ip string) (string, string, error) {
	refreshTokenID := uuid.New()
	accessToken, err := utils.GenerateAccessToken(refreshTokenID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(rawRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	refreshTokenModel := &model.RefreshToken{
		ID:        refreshTokenID,
		UserID:    userID,
		TokenHash: string(hashedRefreshToken),
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: time.Now(),
	}

	if err = s.refreshTokenRepository.Create(refreshTokenModel); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token to database: %w", err)
	}

	return accessToken, rawRefreshToken, nil
}

func (s *AuthServiceImpl) UpdateTokens(refreshTokenID uuid.UUID, rawRefreshToken, userAgent, ip string) (string, string, error) {
	refreshTokenModel, err := s.refreshTokenRepository.GetByID(refreshTokenID)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(refreshTokenModel.TokenHash), []byte(rawRefreshToken)); err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	userID := refreshTokenModel.UserID

	if refreshTokenModel.UserAgent != userAgent {
		if err := s.refreshTokenRepository.MarkAllAsDeactivatedByUserID(userID); err != nil {
			return "", "", fmt.Errorf("failed to deauthorize user: %w", err)
		}
		return "", "", fmt.Errorf("user-agent mismatch. user deauthorized")
	}

	if refreshTokenModel.IP != ip {
		go func() {
			if err := utils.SendWarningToWebhook(userID, refreshTokenModel.IP, ip, userAgent); err != nil {
				log.Printf("failed to send webhook warning: %v", err)
			}
		}()
	}

	if err := s.refreshTokenRepository.MarkAsDeactivated(refreshTokenModel); err != nil {
		return "", "", fmt.Errorf("failed to mark refresh token as used: %w", err)
	}

	newRefreshTokenID := uuid.New()
	newAccessToken, err := utils.GenerateAccessToken(newRefreshTokenID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	newHashedRefreshHash, err := bcrypt.GenerateFromPassword([]byte(newRawRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	newRefreshToken := &model.RefreshToken{
		ID:        newRefreshTokenID,
		UserID:    userID,
		TokenHash: string(newHashedRefreshHash),
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: time.Now(),
	}
	if err := s.refreshTokenRepository.Create(newRefreshToken); err != nil {
		return "", "", fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return newAccessToken, newRawRefreshToken, nil
}

func (s *AuthServiceImpl) DeauthorizeUser(userID uuid.UUID) error {
	if err := s.refreshTokenRepository.MarkAllAsDeactivatedByUserID(userID); err != nil {
		return fmt.Errorf("failed to deauthorize user: %w", err)
	}
	return nil
}

func (s *AuthServiceImpl) GetUserIDByRefreshTokenID(refreshTokenID uuid.UUID) (uuid.UUID, error) {
	refreshTokenModel, err := s.refreshTokenRepository.GetByID(refreshTokenID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("refresh token not found: %w", err)
	}
	return refreshTokenModel.UserID, nil
}

func (s *AuthServiceImpl) IsTokenValid(refreshTokenID uuid.UUID) (bool, error) {
	return s.refreshTokenRepository.IsTokenActive(refreshTokenID)
}

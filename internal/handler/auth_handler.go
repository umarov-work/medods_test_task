package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"medods_test_task/internal/dto"
	"medods_test_task/internal/middleware"
	"medods_test_task/internal/service/intf"
)

type AuthHandler struct {
	authService intf.AuthService
}

func NewAuthHandler(authService intf.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterAuthHandlers(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/create-tokens", h.CreateTokens)
	}

	protected := authGroup.Group("/")
	protected.Use(middleware.AuthMiddleware(h.authService))
	{
		protected.POST("/update-tokens", h.UpdateTokens)
		protected.GET("/deauthorize", h.DeauthorizeUser)
		protected.GET("/me", h.GetUserID)
	}
}

// CreateTokens godoc
// @Summary      Создание access и refresh токенов
// @Description  Генерирует новые токены по userID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "User ID (UUID)"	example(b1506a51-c5a7-45ae-9f2c-4cf700365e46)
// @Success      200      {object}  dto.TokensResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /auth/create-tokens [get]
func (h *AuthHandler) CreateTokens(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "user_id query parameter is required",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid user_id format (must be UUID)",
		})
		return
	}

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	accessToken, refreshToken, err := h.authService.CreateTokens(userID, userAgent, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// UpdateTokens godoc
// @Summary      Обновление access и refresh токенов
// @Description  Обновляет токены по access и refresh токенам
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        updateTokensRequest  body  dto.UpdateTokensRequest  true  "Refresh Token Input"
// @Success      200   {object}  dto.TokensResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /auth/update-tokens [post]
func (h *AuthHandler) UpdateTokens(c *gin.Context) {
	var input dto.UpdateTokensRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	refreshTokenID, err := h.getRefreshTokenIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	accessToken, refreshToken, err := h.authService.UpdateTokens(refreshTokenID, input.RefreshToken, userAgent, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// DeauthorizeUser godoc
// @Summary      Деавторизация пользователя
// @Description  Деактивирует все refresh токены по userID
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /auth/deauthorize [get]
func (h *AuthHandler) DeauthorizeUser(c *gin.Context) {
	refreshTokenID, err := h.getRefreshTokenIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	userID, err := h.authService.GetUserIDByRefreshTokenID(refreshTokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := h.authService.DeauthorizeUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "user deauthorized",
	})
}

// GetUserID godoc
// @Summary      Получение ID пользователя
// @Description  Получает userID пользователя по refreshTokenID из контекста
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UserIDResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *AuthHandler) GetUserID(c *gin.Context) {
	refreshTokenID, err := h.getRefreshTokenIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	userID, err := h.authService.GetUserIDByRefreshTokenID(refreshTokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, dto.UserIDResponse{
		UserID: userID,
	})
}

func (h *AuthHandler) getRefreshTokenIDFromContext(c *gin.Context) (uuid.UUID, error) {
	refreshTokenIDVal, exists := c.Get("refreshTokenID")
	if !exists {
		return uuid.Nil, errors.New("refreshTokenID not found in context")
	}

	refreshTokenID, ok := refreshTokenIDVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid refreshTokenID type")
	}

	return refreshTokenID, nil
}

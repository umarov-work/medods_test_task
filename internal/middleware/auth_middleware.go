package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"medods_test_task/internal/config"
	"medods_test_task/internal/dto"
	service "medods_test_task/internal/service/intf"
)

type AccessTokenClaims struct {
	RefreshTokenID uuid.UUID `json:"refresh_token_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "authorization header is missing",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "authorization header format must be Bearer {token}",
			})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != "HS512" {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.Load().JWTSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*AccessTokenClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "invalid token claims",
			})
			return
		}

		if claims.RefreshTokenID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "invalid token claims: empty refresh_token_id",
			})
			return
		}

		isActive, err := authService.IsTokenValid(claims.RefreshTokenID)
		if err != nil || !isActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "token is no longer valid",
			})
			return
		}

		c.Set("refreshTokenID", claims.RefreshTokenID)

		c.Next()
	}
}

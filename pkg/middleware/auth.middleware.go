package middleware

import (
	"net/http"
	"strings"

	"github.com/Anurag-S1ngh/carbon-backend/pkg/token/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddlewareConfig struct {
	jwtConfig *jwt.JwtConfig
}

func NewAuthMiddlewareConfig(jwtConfig *jwt.JwtConfig) *AuthMiddlewareConfig {
	return &AuthMiddlewareConfig{
		jwtConfig: jwtConfig,
	}
}

func (m *AuthMiddlewareConfig) VerifyAccessToken(c *gin.Context) {
	accessToken, err := c.Cookie("carbon-access-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		c.Abort()
		return
	}

	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		c.Abort()
		return
	}

	userIDStr, err := m.jwtConfig.VerifyToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		c.Abort()
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil || userUUID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		c.Abort()
		return
	}

	c.Set("userUUID", userUUID)
	c.Next()
}

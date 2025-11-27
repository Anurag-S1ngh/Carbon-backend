package http

import (
	"net/http"

	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/http/handler"
	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(authService *service.AuthService) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	authHandler := handler.NewAuthHandler(authService)

	v1 := r.Group("/api/v1")

	v1.POST("/verify-email", authHandler.VerifyEmail)
	v1.POST("/verify-otp", authHandler.VerifyOTP)
	v1.POST("/refresh-access-token", authHandler.RefreshAccessToken)

	return r
}

package http

import (
	"net/http"
	"time"

	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/http/handler"
	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(authService *service.AuthService, redirectURL string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://traccker.anuragcode.me"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	authHandler := handler.NewAuthHandler(authService, redirectURL)

	v1 := r.Group("/api/v1/auth")

	v1.POST("/verify-email", authHandler.VerifyEmail)
	v1.POST("/verify-otp", authHandler.VerifyOTP)
	v1.POST("/refresh-access-token", authHandler.RefreshAccessToken)
	v1.GET("/:provider", authHandler.OAuthProvider)
	v1.GET("/:provider/callback", authHandler.CallbackHandler)

	return r
}

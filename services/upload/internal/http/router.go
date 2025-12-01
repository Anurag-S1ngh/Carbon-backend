package http

import (
	"net/http"

	"github.com/Anurag-S1ngh/carbon-backend/pkg/middleware"
	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/http/handler"
	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(gitService *service.GitService, uploadService *service.UploadService, authMiddlewareConfig *middleware.AuthMiddlewareConfig) *gin.Engine {
	r := gin.Default()

	gitHandler := handler.NewGitHandler(gitService)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	v1 := r.Group("/api/v1")

	v1.GET("/repo/all", authMiddlewareConfig.VerifyAccessToken, gitHandler.ListAllRepos)
	v1.POST("/repo/upload")

	return r
}

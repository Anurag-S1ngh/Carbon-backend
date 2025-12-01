package handler

import (
	"net/http"

	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GitHandler struct {
	gitService *service.GitService
}

func NewGitHandler(gitService *service.GitService) *GitHandler {
	return &GitHandler{
		gitService: gitService,
	}
}

func (h *GitHandler) ListAllRepos(c *gin.Context) {
	userUUID, ok := c.Get("userUUID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in again"})
		return
	}

	repos, err := h.gitService.ListAllRepos(c, userUUID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repos": repos})
}

package handler

import (
	"net/http"
	"strings"

	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	redirectURL string
}

func NewAuthHandler(authService *service.AuthService, redirectURL string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		redirectURL: redirectURL,
	}
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.authService.VerifyEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while sending otp. please try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OTP sent"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.authService.VerifyOTP(c, body.Email, body.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User verified"})
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("carbon-refresh-token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token not found"})
		return
	}

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token not found"})
		return
	}

	err = h.authService.RefreshAccessToken(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sign in again"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "access token refreshed"})
}

func (h *AuthHandler) OAuthProvider(c *gin.Context) {
	h.authService.OAuthProvider(c)
}

func (h *AuthHandler) CallbackHandler(c *gin.Context) {
	err := h.authService.GithubCallbackHandler(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "try again later"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.redirectURL)
}

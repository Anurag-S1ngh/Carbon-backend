package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OTP sent"})
}

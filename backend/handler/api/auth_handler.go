package api

import (
	"net/http"

	"backend/domain/models"
	"backend/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) MockLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Mocking users based on email
	var userID uint
	var role string

	switch req.Email {
	case "admin@finance.com":
		userID = 1
		role = string(models.RoleAdmin)
	case "analyst@finance.com":
		userID = 2
		role = string(models.RoleAnalyst)
	case "viewer@finance.com":
		userID = 3
		role = string(models.RoleViewer)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials. Try admin@, analyst@, or viewer@finance.com"})
		return
	}

	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  role,
	})
}

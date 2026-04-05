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

type LoginRequest struct {
	Email string `json:"email" binding:"required" example:"admin@finance.com"`
}

// MockLogin godoc
// @Summary Mock login for authentication
// @Description Logs in a user mockly (admin@test.com, analyst@test.com, viewer@test.com) and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/auth/login [post]
func (h *AuthHandler) MockLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Mocking users based on email
	var userID uint
	var role string

	switch req.Email {
	case "admin@finance.com", "admin@test.com":
		userID = 1
		role = string(models.RoleAdmin)
	case "analyst@finance.com", "analyst@test.com":
		userID = 2
		role = string(models.RoleAnalyst)
	case "viewer@finance.com", "viewer@test.com":
		userID = 3
		role = string(models.RoleViewer)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials. Try admin@, analyst@, or viewer@test.com"})
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

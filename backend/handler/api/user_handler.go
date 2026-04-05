package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	_ "backend/domain/models"
	"backend/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type CreateUserInput struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	Name  string `json:"name" binding:"required" example:"John Doe"`
}

type UpdateUserStatusInput struct {
	IsActive bool `json:"is_active"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body CreateUserInput true "User payload"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(input.Email, input.Name)
	if err != nil {
		// Detect unique constraint violations (Duplicate Email) to return accurate HTTP status codes natively.
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "User already exists",
				"details": "A user with this email address is already registered in the system.",
			})
			return
		}

		slog.Error("CreateUser Error", "error", err, "email", input.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Fetch a user by ID
// @Description Fetch details of a specific user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary Fetch all users
// @Description Returns a list of all users in the system
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUserStatus godoc
// @Summary Update user status (Active/Inactive)
// @Description Activates or deactivates a user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param status body UpdateUserStatusInput true "Status payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/users/{id}/status [patch]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.userService.UpdateUserStatus(uint(id), input.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

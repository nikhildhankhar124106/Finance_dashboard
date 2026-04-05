package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain/models"
	"backend/handler/api"
	"backend/handler/middleware"
	"backend/infrastructure/db"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.ActivityLog{})
	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userSvc)
	authHandler := api.NewAuthHandler()

	r.POST("/api/v1/auth/login", authHandler.MockLogin)
	
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RequireAuth(db))
	{
		v1.GET("/users", userHandler.GetUsers)
	}

	return r
}

func TestLoginAndAccess(t *testing.T) {
	database := setupTestDB()
	
	// Seed admin user
	database.Create(&models.User{ID: 1, Email: "admin@test.com", Name: "Admin", Role: models.RoleAdmin, IsActive: true})
	
	router := setupTestRouter(database)

	// 1. Test Login
	loginBody := map[string]string{"email": "admin@test.com"}
	body, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"]
	assert.NotEmpty(t, token)

	// 2. Test Protected Resource Access
	req, _ = http.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInactiveUserAccess(t *testing.T) {
	database := setupTestDB()
	database.Create(&models.User{ID: 2, Email: "viewer@test.com", Name: "Viewer", Role: models.RoleViewer, IsActive: false})
	
	router := setupTestRouter(database)

	// Login manually to get token (RequireAuth checks DB every time now)
	reqBody := map[string]string{"email": "viewer@test.com"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"]

	// Access protected route
	req, _ = http.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect Forbidden due to IsActive: false
	assert.Equal(t, http.StatusForbidden, w.Code)
}

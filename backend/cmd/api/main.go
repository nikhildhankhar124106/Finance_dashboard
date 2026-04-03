package main

import (
	"log"
	"net/http"

	"backend/config"
	_ "backend/docs"
	"backend/domain/models"
	"backend/handler/api"
	"backend/handler/middleware"
	"backend/infrastructure/db"
	"backend/pkg/logger"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Finance Dashboard API
// @version 1.0
// @description REST API documentation for the Clean Architecture Golang Finance Dashboard backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

func main() {
	// Initialize Logger
	logger.Init()
	logger.Info.Println("Starting up the backend service...")

	// Load Configuration
	cfg := config.LoadConfig()

	// Connect to Database
	err := db.ConnectPostgres(cfg)
	if err != nil {
		logger.Error.Fatalf("Could not initialize database connection: %v", err)
	}

	// Auto-migrate database models
	DB := db.GetDB()
	err = DB.AutoMigrate(&models.User{}, &models.Transaction{})
	if err != nil {
		logger.Error.Fatalf("Failed to auto migrate database schema: %v", err)
	}

	// Dependency Injection Setup
	userRepo := repository.NewUserRepository(DB)
	userSvc := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userSvc)
	authHandler := api.NewAuthHandler()

	txRepo := repository.NewTransactionRepository(DB)
	txSvc := service.NewTransactionService(txRepo)
	txHandler := api.NewTransactionHandler(txSvc)

	analyticsRepo := repository.NewAnalyticsRepository(DB)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)
	analyticsHandler := api.NewAnalyticsHandler(analyticsSvc)

	// Setup Gin Router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	
	// Apply Central Error Middleware Globally
	r.Use(middleware.ErrorHandler())

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation route
	// To actually serve this, we must run `swag init -g cmd/api/main.go` and import `backend/docs`
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRoutes := r.Group("/api/v1")
	{
		// Public Routes
		apiRoutes.POST("/auth/login", authHandler.MockLogin)

		// Protected Viewer Routes (accessible by Viewer, Analyst, Admin)
		viewerRoutes := apiRoutes.Group("/")
		viewerRoutes.Use(middleware.RequireAuth())
		viewerRoutes.Use(middleware.RequireRole(string(models.RoleViewer), string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			viewerRoutes.GET("/users", userHandler.GetUsers) // Read-only
			viewerRoutes.GET("/transactions", txHandler.GetTransactions)
		}

		// Protected Analyst Routes (accessible by Analyst, Admin)
		analystRoutes := apiRoutes.Group("/")
		analystRoutes.Use(middleware.RequireAuth())
		analystRoutes.Use(middleware.RequireRole(string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			analystRoutes.GET("/summary", analyticsHandler.GetSummary)
			analystRoutes.GET("/category-breakdown", analyticsHandler.GetCategoryBreakdown)
			analystRoutes.GET("/monthly-trends", analyticsHandler.GetMonthlyTrends)
		}

		// Protected Admin Routes (accessible by Admin ONLY)
		adminRoutes := apiRoutes.Group("/")
		adminRoutes.Use(middleware.RequireAuth())
		adminRoutes.Use(middleware.RequireRole(string(models.RoleAdmin)))
		{
			adminRoutes.POST("/users", userHandler.CreateUser) // Full access to create
			adminRoutes.POST("/transactions", txHandler.CreateTransaction)
			adminRoutes.PUT("/transactions/:id", txHandler.UpdateTransaction)
			adminRoutes.DELETE("/transactions/:id", txHandler.DeleteTransaction)
			
			adminRoutes.DELETE("/system/logs", func(c *gin.Context) {
				c.JSON(200, gin.H{"data": "System logs deleted."})
			})
		}
	}

	// Start the server
	logger.Info.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

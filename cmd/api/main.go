package main

import (
	"log"
	"os"
	"ozinse-backend/internal/handler"
	"ozinse-backend/internal/middleware"
	"ozinse-backend/internal/repository"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "ozinse-backend/cmd/api/docs" 
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ozinse API
// @version 1.0
// @description API Server for Ozinse Video Platform Backend
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Initialize a database connection
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. Check database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL on port 5433")

	// 4. Initialize dependencies
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	ageCategoryRepo := repository.NewAgeCategoryRepository(db)
	ageCategoryService := service.NewAgeCategoryService(ageCategoryRepo)
	ageCategoryHandler := handler.NewAgeCategoryHandler(ageCategoryService)

	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	favoriteRepo := repository.NewFavoriteRepository(db)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	// --- NEW: Auth & User Dependencies ---
	// Fetch the secret key from .env, or use a fallback for local development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-super-secret-key"
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	// Pass the service, the secret key, and the token expiration time (e.g., 24 hours)
	authHandler := handler.NewAuthHandler(authService, jwtSecret, 24)
	userHandler := handler.NewUserHandler(userService)

	// 5. Setup a Gin router
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Adding a welcome message
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Ozinse API!")
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "online",
			"project":  "Ozinse API",
			"version":  "1.0.0",
			"database": "connected",
		})
	})

	// 6. Define routes
	api := r.Group("/api")
	{
		// 6.1 Public Routes (No token required)
		api.POST("/auth/login", authHandler.Login) // The counter where users get their token
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/forgot-password", authHandler.ForgotPassword)
		api.POST("/auth/reset-password", authHandler.ResetPassword)

		// Anyone can view data
		api.GET("/categories", categoryHandler.GetAll)
		api.GET("/categories/:id", categoryHandler.GetByID)
		api.GET("/age-categories", ageCategoryHandler.GetAll)
		api.GET("/age-categories/:id", ageCategoryHandler.GetByID)
		api.GET("/genres", genreHandler.GetAll)
		api.GET("/genres/:id", genreHandler.GetByID)
		api.GET("/projects", projectHandler.GetAll)
		api.GET("/projects/:id", projectHandler.GetByID)
		api.GET("/projects/trending", projectHandler.GetTrending)
		api.GET("/projects/:id/similar", projectHandler.GetSimilar)

		// 6.2 Protected Routes (Require a valid JWT token)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{

			protected.GET("/users/:id", userHandler.GetByID)
			protected.PUT("/users/:id", userHandler.Update)

			protected.POST("/users/favorites", favoriteHandler.AddFavorite)
			protected.DELETE("/users/favorites/:project_id", favoriteHandler.RemoveFavorite)

			adminOnly := protected.Group("/")
			adminOnly.Use(middleware.AdminOnly())
			{
				// only users with role_id = 2 (admin) can access these routes
				adminOnly.POST("/categories", categoryHandler.Create)
				adminOnly.PUT("/categories/:id", categoryHandler.Update)
				adminOnly.DELETE("/categories/:id", categoryHandler.Delete)

				adminOnly.POST("/age-categories", ageCategoryHandler.Create)
				adminOnly.PUT("/age-categories/:id", ageCategoryHandler.Update)
				adminOnly.DELETE("/age-categories/:id", ageCategoryHandler.Delete)

				adminOnly.POST("/genres", genreHandler.Create)
				adminOnly.PUT("/genres/:id", genreHandler.Update)
				adminOnly.DELETE("/genres/:id", genreHandler.Delete)

				adminOnly.POST("/projects", projectHandler.Create)
				adminOnly.PUT("/projects/:id", projectHandler.Update)
				adminOnly.DELETE("/projects/:id", projectHandler.Delete)

				adminOnly.POST("/roles", roleHandler.CreateRole)
				adminOnly.GET("/roles", roleHandler.GetAllRoles)
				adminOnly.GET("/roles/:id", roleHandler.GetRoleByID)
				adminOnly.PUT("/roles/:id", roleHandler.UpdateRole)
				adminOnly.DELETE("/roles/:id", roleHandler.DeleteRole)

				adminOnly.GET("/users", userHandler.GetAll)
				adminOnly.DELETE("/users/:id", userHandler.Delete)
				
				adminOnly.POST("/main-page-entries", projectHandler.CreateMainPageEntry)
				adminOnly.GET("/main-page-entries", projectHandler.GetMainPageEntries)
				adminOnly.DELETE("/main-page-entries/:id", projectHandler.DeleteMainPageEntry)
				adminOnly.PUT("/main-page-entries/:id", projectHandler.UpdateMainPageEntry)
				adminOnly.GET("/main-page-entries/:id", projectHandler.GetByIDForMainPage)
				

			}
		}
	}

	// 7. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

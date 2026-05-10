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
)

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

	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// --- NEW: Auth & User Dependencies ---
	// Fetch the secret key from .env, or use a fallback for local development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-super-secret-key"
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	// Pass the service, the secret key, and the token expiration time (e.g., 24 hours)
	authHandler := handler.NewAuthHandler(userService, jwtSecret, 24)

	// 5. Setup a Gin router
	r := gin.Default()

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

		// Anyone can view data
		api.GET("/categories", categoryHandler.GetAll)
		api.GET("/categories/:id", categoryHandler.GetByID)
		api.GET("/genres", genreHandler.GetAll)
		api.GET("/genres/:id", genreHandler.GetByID)
		api.GET("/projects", projectHandler.GetAll)
		api.GET("/projects/:id", projectHandler.GetByID)

		// 6.2 Protected Routes (Require a valid JWT token)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
	
			adminOnly := protected.Group("/")
			adminOnly.Use(middleware.AdminOnly())
			{
				// only users with role_id = 2 (admin) can access these routes
				adminOnly.POST("/categories", categoryHandler.Create)
				adminOnly.PUT("/categories/:id", categoryHandler.Update)
				adminOnly.DELETE("/categories/:id", categoryHandler.Delete)

				adminOnly.POST("/genres", genreHandler.Create)
				adminOnly.PUT("/genres/:id", genreHandler.Update)
				adminOnly.DELETE("/genres/:id", genreHandler.Delete)

				adminOnly.POST("/projects", projectHandler.Create)
				adminOnly.PUT("/projects/:id", projectHandler.Update)
				adminOnly.DELETE("/projects/:id", projectHandler.Delete)
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
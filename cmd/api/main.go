package main

import (
	"log"
	"os"
	"ozinse-backend/internal/handler"
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
	
	projectRepo:= repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

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
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		genres := api.Group("/genres")
		{
			genres.GET("", genreHandler.GetAll)
			genres.GET("/:id", genreHandler.GetByID)
			genres.POST("", genreHandler.Create)
			genres.PUT("/:id", genreHandler.Update)
			genres.DELETE("/:id", genreHandler.Delete)
		}

		projects := api.Group("/projects")
		{
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.POST("", projectHandler.Create)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
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
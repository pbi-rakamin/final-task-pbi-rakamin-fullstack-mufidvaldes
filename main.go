package main

import (
	"example.com/g-auth/database"
	"example.com/g-auth/middleware"
	"example.com/g-auth/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Initialize database connection
	database.InitDB()
	// Initialize validator middleware
	middleware.SetupValidation()

	// Apply middleware for validation
	router.Use(func(c *gin.Context) {
		c.Set("validate", middleware.ValidateInput)
		c.Next()
	})

	// Apply middleware for authentication
	router.Use(middleware.AuthMiddleware())

	routes.SetupUserRoutes(router)
	routes.SetupPhotoRoutes(router)

	router.Run(":8080")
}

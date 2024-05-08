package routes

import (
	"example.com/g-auth/controllers"
	"github.com/gin-gonic/gin"
)

// User Endpoints
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", controllers.RegisterUser)
		userGroup.POST("/login", controllers.LoginUser)
		userGroup.PUT("/:userId", controllers.UpdateUser)
		userGroup.DELETE("/:userId", controllers.DeleteUser)

	}
}

// SetupPhotoRoutes mengatur rute-rute untuk foto.
func SetupPhotoRoutes(router *gin.Engine) {
	photoGroup := router.Group("/photos")
	{
		// Photo Endpoints
		photoGroup.POST("/", controllers.CreatePhoto)
		photoGroup.GET("/", controllers.GetPhotos)
		photoGroup.PUT("/:photoId", controllers.UpdatePhoto)
		photoGroup.DELETE("/:photoId", controllers.DeletePhoto)

	}
}

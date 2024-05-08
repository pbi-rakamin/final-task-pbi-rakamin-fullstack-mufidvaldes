package controllers

import (
	"example.com/g-auth/database"
	"example.com/g-auth/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePhoto creates a new photo
func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Mendapatkan user_id dari token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Missing JWT claims"})
		return
	}
	userID := claims.(jwt.MapClaims)["userID"].(float64)
	// Mengatur user_id foto sesuai dengan user yang sedang login
	photo.UserID = uint(userID)

	if err := database.DB.Create(&photo).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create photo"})
		return
	}

	c.JSON(200, gin.H{"data": photo})
}

// GetPhotos retrieves all photos
func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	if err := database.DB.Find(&photos).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch photos"})
		return
	}

	c.JSON(200, gin.H{"data": photos})
}

// UpdatePhoto updates a photo's information
func UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")
	var updatedPhoto models.Photo
	if err := c.ShouldBindJSON(&updatedPhoto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var photo models.Photo
	if err := database.DB.First(&photo, photoID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to update photo"})
		return
	}
	// Mendapatkan user_id dari token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Missing JWT claims"})
		return
	}
	userID := claims.(jwt.MapClaims)["userID"].(float64)
	if photo.UserID != uint(userID) {
		c.JSON(403, gin.H{"error": "You are not authorized to update this photo"})
		return
	}
	// Update photo fields
	photo.Title = updatedPhoto.Title
	photo.Caption = updatedPhoto.Caption
	photo.PhotoURL = updatedPhoto.PhotoURL

	if err := database.DB.Save(&photo).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update photo"})
		return
	}

	c.JSON(200, gin.H{"data": photo})
}

// DeletePhoto deletes a photo
func DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := database.DB.First(&photo, photoID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete photo"})
		return
	}
	// Mendapatkan user_id dari token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Missing JWT claims"})
		return
	}
	userID := claims.(jwt.MapClaims)["userID"].(float64)
	if photo.UserID != uint(userID) {
		c.JSON(403, gin.H{"error": "You are not authorized to update this photo"})
		return
	}

	if err := database.DB.Delete(&photo).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(200, gin.H{"message": "Photo deleted successfully"})
}

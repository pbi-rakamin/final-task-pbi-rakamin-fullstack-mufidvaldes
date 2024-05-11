package controllers

import (
	"net/http"
	"strconv"
	"time"

	"example.com/g-auth/database"
	"example.com/g-auth/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// RegisterUser registers a new user
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Validate input
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(200, gin.H{"data": user})
}

// LoginUser authenticates a user and generates JWT token
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to login"})
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(existingUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

// UpdateUser updates a user's information
func UpdateUser(c *gin.Context) {
	userID := c.Param("userId")

	// Mendapatkan user_id dari token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Missing JWT claims"})
		return
	}
	authenticatedUserID := uint(claims.(jwt.MapClaims)["userID"].(float64))

	// Memeriksa apakah pengguna yang sedang mengakses profil adalah pengguna yang terautentikasi
	if userID != strconv.Itoa(int(authenticatedUserID)) {
		c.JSON(403, gin.H{"error": "You are not allowed to delete another user's profile"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	// Update user fields
	user.Username = updatedUser.Username
	user.Email = updatedUser.Email
	// Jika password juga diperbarui, hash password baru
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	// Generate new JWT token
	newToken, err := generateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(200, gin.H{"data": user, "token": newToken})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	// Mendapatkan user_id dari token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Missing JWT claims"})
		return
	}
	authenticatedUserID := uint(claims.(jwt.MapClaims)["userID"].(float64))

	// Memeriksa apakah pengguna yang sedang mengakses profil adalah pengguna yang terautentikasi
	if userID != strconv.Itoa(int(authenticatedUserID)) {
		c.JSON(403, gin.H{"error": "You are not allowed to delete another user's profile"})
		return
	}
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// generateToken generates JWT token for user authentication
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response
	return token.SignedString([]byte("secret-key"))
}

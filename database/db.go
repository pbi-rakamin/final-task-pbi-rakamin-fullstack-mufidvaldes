package database

import (
	"example.com/g-auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:@(127.0.0.1:3306)/db-g-base?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
	// Auto Migrate the Models
	db.AutoMigrate(&models.User{}, &models.Photo{})
}

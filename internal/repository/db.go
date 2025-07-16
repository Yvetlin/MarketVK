package repository

import (
	"MarketVK/internal/model"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() error {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	if err := DB.AutoMigrate(&model.User{}, &model.Ad{}); err != nil {
		return err
	}
	return nil
}

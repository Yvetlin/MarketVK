package db

import (
	"MarketVK/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestDB() error {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	if err := DB.AutoMigrate(&model.User{}, &model.Ad{}); err != nil {
		return err
	}
	return nil
}

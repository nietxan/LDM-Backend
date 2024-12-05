package web

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error

	DB, err = gorm.Open(sqlite.Open("lmd.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := DB.AutoMigrate(&User{}, &Order{}); err != nil {
		return err
	}
	
	return nil
}

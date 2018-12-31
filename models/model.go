package models

import (
	"../utils"

	"github.com/jinzhu/gorm"
)

const (
	AdminName     = "Admin"
	AdminLogin    = "admin"
	AdminPassword = "admin123"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Group{}, &Task{}, &Tag{})

	// Create User Admin
	hash, err := utils.HashPassword(AdminPassword)
	if err != nil {
		panic(err)
	}

	db.Create(&User{Name: AdminName, Login: AdminLogin, Password: hash, IsAdmin: true})
}

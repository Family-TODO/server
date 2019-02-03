package models

import (
	"../utils"

	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

const (
	adminName     = "Admin"
	adminLogin    = "admin"
	adminPassword = "admin123"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Group{}, &Task{})

	// Create User Admin
	hash, err := utils.HashPassword(adminPassword)
	if err != nil {
		panic(err)
	}

	db.Create(&User{Name: adminName, Login: adminLogin, Password: hash, IsAdmin: true})
}

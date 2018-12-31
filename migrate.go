package main

import (
	"./config"
	"./models"
)

const (
	AdminName     = "Admin"
	AdminLogin    = "admin"
	AdminPassword = "admin123"
)

func main() {
	db, session, _ := config.Init()
	defer db.Close()
	defer session.Close()

	// Migrate the schema
	models.Migrate(db)

	// Create User Admin
	hash, err := models.HashPassword(AdminPassword)
	if err != nil {
		panic(err)
	}

	db.Create(&models.User{Name: AdminName, Login: AdminLogin, Password: hash, IsAdmin: true})
}

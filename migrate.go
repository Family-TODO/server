package main

import (
	"./models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

const (
	AdminName     = "Admin"
	AdminLogin    = "admin"
	AdminPassword = "admin123"
)

func main() {
	/* - Import Environment - */
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	/* - Connect to Database - */
	db, err := gorm.Open("sqlite3", os.Getenv("DATABASE_PATH"))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Migrate the schema
	models.Migrate(db)

	// Create User Admin
	hash, err := models.HashPassword(AdminPassword)
	if err != nil {
		panic(err)
	}

	db.Create(&models.User{Name: AdminName, Login: AdminLogin, Password: hash, IsAdmin: true})
}

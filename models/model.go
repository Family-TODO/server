package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const HashCost = 14

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &UserSession{})

	db.Model(&UserSession{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
}

// Password Hashing.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)
	return string(bytes), err
}

// Validate Password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

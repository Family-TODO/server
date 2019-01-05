package utils

import "golang.org/x/crypto/bcrypt"

const HashCost = 14

// Password Hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)
	return string(bytes), err
}

// Validate Password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

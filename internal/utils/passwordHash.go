package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bcryptCost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

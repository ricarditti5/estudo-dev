package main

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bcryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptedPass), nil
}

func CheckPasswordHash(password, hash string) bool {
	checkPassword := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return checkPassword == nil
}

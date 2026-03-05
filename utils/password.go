package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HasPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic("failed to hash password: " + err.Error())
	}
	return string(hashed)
}

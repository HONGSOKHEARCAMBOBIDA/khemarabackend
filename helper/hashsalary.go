package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HasSalary(salary string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(salary), bcrypt.MinCost)
	if err != nil {
		panic("failed to hash password: " + err.Error())
	}
	return string(hashed)
}

func MaskSalary(salary string) string {
	if len(salary) <= 3 {
		return "***"
	}
	return "***" + salary[len(salary)-3:]
}

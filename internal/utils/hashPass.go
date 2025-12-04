package utils

import "golang.org/x/crypto/bcrypt"

func HashPass(pass string) (string, error) {
	cost := 12
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(hashed), err
}

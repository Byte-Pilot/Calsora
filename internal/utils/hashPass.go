package utils

import "golang.org/x/crypto/bcrypt"

func HashPass(pass string) (string, error) {
	cost := 12
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(hashed), err
}

func CheckPass(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

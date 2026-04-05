package validator

import (
	"errors"
	"regexp"
)

var errPassLength = errors.New("incorrect password length")
var errPassInvalid = errors.New("invalid password ")

func ValidatePass(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errPassLength
	}
	passOk := regexp.MustCompile(`^[a-zA-Z0-9!?$%()_.,-]+$`)
	if !passOk.MatchString(password) {
		return errPassInvalid
	}
	return nil
}

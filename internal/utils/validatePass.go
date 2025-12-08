package utils

import (
	"errors"
	"regexp"
)

var errPassTooShort = errors.New("password too short")
var errPassInvalid = errors.New("password invalid")

func ValidatePass(password string) error {
	if len(password) < 8 {
		return errPassTooShort
	}
	passOk := regexp.MustCompile(`^[a-zA-Z0-9!?()-_.,]+$`)
	if !passOk.MatchString(password) {
		return errPassInvalid
	}
	return nil
}

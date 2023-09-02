package utils

import (
	"net/mail"
	"regexp"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	validLength := len(password) >= 8 && len(password) <= 30
	numbers := regexp.MustCompile("[0-9]+").FindString(password)
	lower := regexp.MustCompile("[a-z]+").FindString(password)
	caps := regexp.MustCompile("[A-Z]+").FindString(password)
	return validLength && numbers != "" && lower != "" && caps != ""
}

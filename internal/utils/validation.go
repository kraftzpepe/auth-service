package utils

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	if len(email) < 3 || !regexp.MustCompile(`^[^@]+@[^@]+\.(com)$`).MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

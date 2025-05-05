package auth

import (
	"errors"
	"regexp"
)

var PasswordValidationError = errors.New("password must be between 6 to 16 characters, contain at least one uppercase letter, one lowercase letter, and one special character")

func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 16 {
		return PasswordValidationError
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasSpecial := regexp.MustCompile(`[(){}\[\]<>+\-*/?,.:;"'_\|~` + "`" + `!@#$%^&=]`).MatchString

	if !hasUpper(password) || !hasLower(password) || !hasSpecial(password) {
		return PasswordValidationError
	}

	return nil
}

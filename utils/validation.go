package utils

import (
	"errors"
	"github.com/AdiInfiniteLoop/Authora/models"
	"regexp"
)

type ValidationResult struct {
	isValid bool
	Error   error
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewValidationResult(isValid bool, err error) ValidationResult {
	return ValidationResult{
		isValid: isValid,
		Error:   err,
	}
}

func ValidateEmail(email string) ValidationResult {
	if email == "" {
		return NewValidationResult(false, errors.New("empty email is not allowed"))
	}
	if !emailRegex.MatchString(email) {
		return NewValidationResult(false, errors.New("invalid Email"))
	}
	return NewValidationResult(true, nil)
}

func ValidatePassword(password string) ValidationResult {
	if len(password) < 6 {
		return NewValidationResult(false, errors.New("password must be minimum 6 characters"))
	}
	return NewValidationResult(true, nil)
}

func ValidationOfUser(user models.User) []string {
	var _errors []string
	if v := ValidateEmail(user.Email); !v.isValid {
		_errors = append(_errors, v.Error.Error())
	}
	if v := ValidatePassword(user.Password); !v.isValid {
		_errors = append(_errors, v.Error.Error())
	}
	return _errors
}

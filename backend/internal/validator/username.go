package validator

import (
	"errors"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
)

var (
	ErrUsernameTooShort      = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong       = errors.New("username must be at most 32 characters")
	ErrUsernameInvalidFormat = errors.New("username must contain only lowercase letters, numbers, and hyphens; must start and end with alphanumeric")
)

// ValidateUsername validates a username against K8s naming constraints.
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	if len(username) > 32 {
		return ErrUsernameTooLong
	}
	if !usernameRegex.MatchString(username) {
		return ErrUsernameInvalidFormat
	}
	return nil
}

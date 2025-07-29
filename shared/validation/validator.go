package validation

import (
	"errors"
	"regexp"
)

// ValidateUUID returns error if input is not a valid UUID
func ValidateUUID(uuid string) error {
	re := regexp.MustCompile(`^[a-fA-F0-9]{8}\-[a-fA-F0-9]{4}\-[1-5][a-fA-F0-9]{3}\-[89abAB][a-fA-F0-9]{3}\-[a-fA-F0-9]{12}$`)
	if !re.MatchString(uuid) {
		return errors.New("invalid UUID format")
	}
	return nil
}

// ValidateEmail format (basic check)
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

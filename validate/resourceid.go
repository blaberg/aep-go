package validate

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"github.com/google/uuid"
)

var validMiddleChars = regexp.MustCompile(`^[a-z0-9-]+$`)

// ResourceID validates that a resource ID follows the AEP design specification.
func ResourceID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("resource ID can not be empty")
	}
	if len(id) > 63 {
		return fmt.Errorf("resource ID SHOULD not be longer than 63 characters")
	}
	if !unicode.IsLower(rune(id[0])) {
		return errors.New("resource ID must start with a lowercase letter")
	}
	if id[len(id)-1] == '-' {
		return errors.New("resource ID cannot end with a hyphen")
	}
	if !validMiddleChars.MatchString(id) {
		return errors.New("resource ID can only contain lowercase letters, numbers, and hyphens")
	}
	if err := uuid.Validate(id); err == nil {
		return errors.New("resource ID can not ba a UUID")
	}
	return nil
}

package customvalidator

import (
	"errors"
	"regexp"
)

func ValidateBdPhoneNumber(phone string) error {
	// Define the regex pattern
	pattern := `^(01)[3-9]{1}[0-9]{8}$`

	// Compile the regex
	re := regexp.MustCompile(pattern)

	// Validate the phone number
	if !re.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

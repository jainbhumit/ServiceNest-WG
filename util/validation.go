package util

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"
)

// ValidateEmail checks if the email is valid.
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

// ValidatePassword checks if the password meets security requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol") // Password is too short
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	// Return true if all criteria are met
	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol") // Password is too short

	}

	return nil
}

// ValidatePhoneNumber checks if the phone number is valid.
func ValidatePhoneNumber(phone string) error {

	// Check if the phone number is numeric
	if _, err := strconv.Atoi(phone); err != nil {
		return errors.New("invalid phone number")

	}

	// Check the length of the phone number
	if len(phone) < 10 || len(phone) > 10 { // Assuming a range of valid lengths
		return errors.New("invalid phone number")

	}

	return nil

}

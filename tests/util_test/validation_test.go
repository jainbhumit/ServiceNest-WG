package util_test

import (
	"errors"
	"serviceNest/util"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected error
	}{
		{"test@example.com", nil},
		{"invalid-email", errors.New("invalid email address")},
		{"@example.com", errors.New("invalid email address")},
		{"test@.com", errors.New("invalid email address")},
		{"test@example", errors.New("invalid email address")},
	}

	for _, test := range tests {
		result := util.ValidateEmail(test.email)
		if result != nil && result.Error() != test.expected.Error() {
			t.Errorf("ValidateEmail(%q) = %v; want %v", test.email, result, test.expected)
		} else if result == nil && test.expected != nil {
			t.Errorf("ValidateEmail(%q) = nil; want %v", test.email, test.expected)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected error
	}{
		{"Valid1Password!", nil},
		{"short1!", errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol")},
		{"NoSpecialChar1", errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol")},
		{"noNumber!", errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol")},
		{"12345678", errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol")},
	}

	for _, test := range tests {
		result := util.ValidatePassword(test.password)
		if result != nil && result.Error() != test.expected.Error() {
			t.Errorf("ValidatePassword(%q) = %v; want %v", test.password, result, test.expected)
		} else if result == nil && test.expected != nil {
			t.Errorf("ValidatePassword(%q) = nil; want %v", test.password, test.expected)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		phone    string
		expected error
	}{
		{"1234567890", nil},
		{"12345", errors.New("invalid phone number")},
		{"abcdefghij", errors.New("invalid phone number")},
		{"12345678901", errors.New("invalid phone number")},
		{"123-456-7890", errors.New("invalid phone number")},
	}

	for _, test := range tests {
		result := util.ValidatePhoneNumber(test.phone)
		if result != nil && result.Error() != test.expected.Error() {
			t.Errorf("ValidatePhoneNumber(%q) = %v; want %v", test.phone, result, test.expected)
		} else if result == nil && test.expected != nil {
			t.Errorf("ValidatePhoneNumber(%q) = nil; want %v", test.phone, test.expected)
		}
	}
}

package util_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
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

func TestApplyPagination(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	t.Run("Test valid pagination within bounds", func(t *testing.T) {
		// Apply pagination with limit 3, offset 2
		result := util.ApplyPagination(data, 3, 2)

		expected := []int{3, 4, 5}
		assert.Equal(t, expected, result, "expected a slice of [3, 4, 5]")
	})

	t.Run("Test pagination beyond data length", func(t *testing.T) {
		// Apply pagination with a large offset
		result := util.ApplyPagination(data, 5, 20)

		expected := []int{}
		assert.Equal(t, expected, result, "expected an empty slice as offset is out of bounds")
	})

	t.Run("Test pagination with end exceeding slice length", func(t *testing.T) {
		// Apply pagination with limit 5, offset 8 (beyond the bounds)
		result := util.ApplyPagination(data, 5, 8)

		expected := []int{9, 10} // The last two items should be returned
		assert.Equal(t, expected, result, "expected a slice of [9, 10] when end exceeds the data length")
	})

	t.Run("Test pagination with limit exceeding slice length", func(t *testing.T) {
		// Apply pagination with large limit
		result := util.ApplyPagination(data, 20, 0)

		expected := data // All items should be returned
		assert.Equal(t, expected, result, "expected the whole slice to be returned")
	})
}

func TestGetPaginationParams(t *testing.T) {
	t.Run("Test with valid limit and offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?limit=5&offset=2", nil)

		limit, offset := util.GetPaginationParams(req)

		assert.Equal(t, 5, limit, "expected limit to be 5")
		assert.Equal(t, 2, offset, "expected offset to be 2")
	})

	t.Run("Test with default values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		limit, offset := util.GetPaginationParams(req)

		assert.Equal(t, 10, limit, "expected default limit to be 10")
		assert.Equal(t, 0, offset, "expected default offset to be 0")
	})

	t.Run("Test with invalid limit and offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?limit=invalid&offset=invalid", nil)

		limit, offset := util.GetPaginationParams(req)

		assert.Equal(t, 10, limit, "expected default limit of 10 when limit is invalid")
		assert.Equal(t, 0, offset, "expected default offset of 0 when offset is invalid")
	})

	t.Run("Test with negative offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?offset=-5", nil)

		limit, offset := util.GetPaginationParams(req)

		assert.Equal(t, 10, limit, "expected default limit of 10")
		assert.Equal(t, 0, offset, "expected default offset of 0 when offset is negative")
	})

	t.Run("Test with negative limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?limit=-5", nil)

		limit, offset := util.GetPaginationParams(req)

		assert.Equal(t, 10, limit, "expected default limit of 10 when limit is negative")
		assert.Equal(t, 0, offset, "expected default offset of 0")
	})
}

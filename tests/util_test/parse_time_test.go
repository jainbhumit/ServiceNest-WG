package util_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"serviceNest/util"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name      string
		input     []uint8
		expected  time.Time
		expectErr bool
	}{
		{
			name:      "Valid time string",
			input:     []uint8("2023-09-08 14:30:00"),
			expected:  time.Date(2023, 9, 8, 14, 30, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			name:      "Invalid time string",
			input:     []uint8("invalid time string"),
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name:      "Empty time string",
			input:     []uint8(""),
			expected:  time.Time{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTime, err := util.ParseTime(tt.input)

			// If we expect an error, assert that the error is not nil
			if tt.expectErr {
				assert.Error(t, err, fmt.Sprintf("expected error but got none for input %v", tt.input))
			} else {
				assert.NoError(t, err, fmt.Sprintf("did not expect error but got %v", err))
				assert.Equal(t, tt.expected, parsedTime, fmt.Sprintf("expected %v, got %v", tt.expected, parsedTime))
			}
		})
	}
}

//	func TestParseTimeAll(t *testing.T) {
//		tests := []struct {
//			name      string
//			input     []uint8
//			expected  time.Time
//			expectErr bool
//		}{
//			{
//				name:      "Valid RFC3339 time string",
//				input:     []uint8("2023-09-08 14:30:00"),
//				expected:  time.Date(2023, 9, 8, 14, 30, 0, 0, time.UTC),
//				expectErr: false,
//			},
//			{
//				name:      "Valid default format time string",
//				input:     []uint8("2023-09-08 14:30:00"),
//				expected:  time.Date(2023, 9, 8, 14, 30, 0, 0, time.UTC),
//				expectErr: false,
//			},
//			{
//				name:      "Invalid time string",
//				input:     []uint8("invalid time string"),
//				expected:  time.Time{}, // zero value of time.Time
//				expectErr: true,
//			},
//			{
//				name:      "Empty time string",
//				input:     []uint8(""),
//				expected:  time.Time{}, // zero value of time.Time
//				expectErr: true,
//			},
//			{
//				name:      "Time string with invalid format but similar to default format",
//				input:     []uint8("2023-09-08 14:30"),
//				expected:  time.Time{}, // zero value of time.Time
//				expectErr: true,
//			},
//			{
//				name:      "Time string with additional characters",
//				input:     []uint8("2023-09-08 14:30:00 extra"),
//				expected:  time.Time{}, // zero value of time.Time
//				expectErr: true,
//			},
//		}
//
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				parsedTime, err := util.ParseTime(tt.input)
//
//				// Check if an error was expected
//				if tt.expectErr {
//					assert.Error(t, err, fmt.Sprintf("expected error but got none for input %v", tt.input))
//					assert.Equal(t, tt.expected, parsedTime, fmt.Sprintf("expected %v, got %v", tt.expected, parsedTime))
//				} else {
//					assert.NoError(t, err, fmt.Sprintf("did not expect error but got %v", err))
//					assert.True(t, tt.expected.Equal(parsedTime), fmt.Sprintf("expected %v, got %v", tt.expected, parsedTime))
//				}
//			})
//		}
//	}
func TestParseTimeAll(t *testing.T) {
	tests := []struct {
		name      string
		input     []uint8
		expected  time.Time
		expectErr bool
	}{
		{
			name:      "Valid RFC3339 time string",
			input:     []uint8("2023-09-08T14:30:00Z"),
			expected:  time.Date(2023, 9, 8, 14, 30, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			name:      "Valid default format time string",
			input:     []uint8("2023-09-08 14:30:00"),
			expected:  time.Date(2023, 9, 8, 14, 30, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			name:      "Valid time string with leading/trailing spaces",
			input:     []uint8(" 2023-09-08 14:30:00 "),
			expected:  time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			expectErr: true,
		},
		{
			name:      "Invalid time string",
			input:     []uint8("invalid time string"),
			expected:  time.Time{}, // zero value of time.Time
			expectErr: true,
		},
		{
			name:      "Empty time string",
			input:     []uint8(""),
			expected:  time.Time{}, // zero value of time.Time
			expectErr: true,
		},
		{
			name:      "Malformed time string",
			input:     []uint8("2023-09-08T14:30:00"),
			expected:  time.Time{}, // zero value of time.Time
			expectErr: true,
		},
		{
			name:      "Time string with additional characters",
			input:     []uint8("2023-09-08 14:30:00 extra"),
			expected:  time.Time{}, // zero value of time.Time
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTime, err := util.ParseTime(tt.input)

			// Check if an error was expected
			if tt.expectErr {
				assert.Error(t, err, fmt.Sprintf("expected error but got none for input %v", tt.input))
				assert.Equal(t, tt.expected, parsedTime, fmt.Sprintf("expected %v, got %v", tt.expected, parsedTime))
			} else {
				assert.NoError(t, err, fmt.Sprintf("did not expect error but got %v", err))
				assert.True(t, tt.expected.Equal(parsedTime), fmt.Sprintf("expected %v, got %v", tt.expected, parsedTime))
			}
		})
	}
}

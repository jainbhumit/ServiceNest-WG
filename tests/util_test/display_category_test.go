package util_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/util"
	"testing"
)

func TestDisplayCategory(t *testing.T) {
	// Mock data
	mockCategories := []model.Category{
		{Name: "Category1", Description: "Description1"},
		{Name: "Category2", Description: "Description2"},
	}

	// Convert mock data to JSON
	mockData, err := json.Marshal(mockCategories)
	if err != nil {
		t.Fatalf("Failed to marshal mock data: %v", err)
	}

	// Override the global Print function to capture output
	var buf bytes.Buffer
	originalPrint := util.Print
	util.Print = func(format string, args ...interface{}) (n int, err error) {
		return fmt.Fprintf(&buf, format, args...)
	}
	defer func() { util.Print = originalPrint }() // Restore original Print after test

	// Mock ioutil.ReadFile to return the mock data
	oldReadFile := util.ReadFile
	defer func() { util.ReadFile = oldReadFile }() // Restore original ReadFile after test
	util.ReadFile = func(filename string) ([]byte, error) {
		return mockData, nil
	}

	// Call the function to test
	util.DisplayCategory()

	// Define the expected output
	expectedOutput := `1 Name : Category1 Description : Description1
2 Name : Category2 Description : Description2
`

	// Verify the output
	assert.Equal(t, expectedOutput, buf.String())
}
func TestDisplayCategoryAll(t *testing.T) {
	tests := []struct {
		name           string
		mockData       []model.Category
		mockReadFile   func(filename string) ([]byte, error)
		expectedOutput string
		expectError    bool
	}{
		{
			name: "Successful Read and Unmarshal",
			mockData: []model.Category{
				{Name: "Category1", Description: "Description1"},
				{Name: "Category2", Description: "Description2"},
			},
			mockReadFile: func(filename string) ([]byte, error) {
				return json.Marshal([]model.Category{
					{Name: "Category1", Description: "Description1"},
					{Name: "Category2", Description: "Description2"},
				})
			},
			expectedOutput: `1 Name : Category1 Description : Description1
2 Name : Category2 Description : Description2
`,
			expectError: false,
		},

		{
			name:           "Empty Category List",
			mockData:       []model.Category{},
			mockReadFile:   func(filename string) ([]byte, error) { return json.Marshal([]model.Category{}) },
			expectedOutput: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the global Print function to capture output
			var buf bytes.Buffer
			originalPrint := util.Print
			util.Print = func(format string, args ...interface{}) (n int, err error) {
				return fmt.Fprintf(&buf, format, args...)
			}
			defer func() { util.Print = originalPrint }() // Restore original Print after test

			// Mock ioutil.ReadFile to return the mock data
			originalReadFile := util.ReadFile
			util.ReadFile = tt.mockReadFile
			defer func() { util.ReadFile = originalReadFile }() // Restore original ReadFile after test

			// Call the function to test
			util.DisplayCategory()

			// Verify the output
			if tt.expectError {
				assert.Contains(t, buf.String(), tt.expectedOutput)
			} else {
				assert.Equal(t, tt.expectedOutput, buf.String())
			}
		})
	}
}

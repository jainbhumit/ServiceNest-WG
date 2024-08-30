package util_test

import (
	"serviceNest/util"
	"strconv"
	"testing"
)

func TestGenerateUniqueID(t *testing.T) {
	// Generate a unique ID
	id := util.GenerateUniqueID()

	// Check if the ID is not empty
	if id == "" {
		t.Errorf("GenerateUniqueID() = %q; want non-empty string", id)
	}

	// Check if the ID is a valid integer within the expected range
	intID, err := strconv.Atoi(id)
	if err != nil {
		t.Errorf("GenerateUniqueID() = %q; want valid integer", id)
	}

	if intID < 0 || intID >= 10000 {
		t.Errorf("GenerateUniqueID() = %d; want integer between 0 and 9999", intID)
	}
}

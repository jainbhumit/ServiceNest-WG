package util

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

func GenerateUniqueID() string {
	return fmt.Sprintf("%d", rand.Intn(10000))
}

func GenerateUUID() string {
	return uuid.New().String()
}

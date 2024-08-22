package util

import (
	"fmt"
	"math/rand"
)

func GenerateUniqueID() string {
	return fmt.Sprintf("%d", rand.Intn(10000))
}

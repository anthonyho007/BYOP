package server

import (
	"fmt"
	"math/rand"
)

func generateId() string {
	id := fmt.Sprintf("%d", 100+rand.Intn(999999))
	return id
}

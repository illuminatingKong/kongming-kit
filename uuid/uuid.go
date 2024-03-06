package uuid

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ID() string {
	return uuid.New().String()
}

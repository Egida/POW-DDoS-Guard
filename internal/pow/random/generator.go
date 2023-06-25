package random

import (
	"math/rand"
	"time"
)

// Generator used for different types values
type Generator interface {
	UInt64() uint64
}

type stdGenerator struct {
	seed *rand.Rand
}

func NewGenerator() Generator {
	return &stdGenerator{
		seed: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

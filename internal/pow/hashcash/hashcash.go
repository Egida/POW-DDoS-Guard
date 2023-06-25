package hashcash

import (
	"hash"
	"math/big"

	"github.com/Goganad/pow_ddos_guard/internal/pow/random"
)

type Config struct {
	Difficulty int
	HashFunc   func() hash.Hash
}

type Hashcash struct {
	hashFunc        func() hash.Hash
	targetBits      *big.Int
	randomGenerator random.Generator
}

type Performer struct {
	hashFunc func() hash.Hash
}

func New(cfg Config) *Hashcash {
	hashBitsCount := cfg.HashFunc().Size() * 8
	// shift to locate first difficulty zeros
	// all bits count minus first difficulty bits count
	targetBitShift := uint(hashBitsCount - cfg.Difficulty)

	// all bits are 1
	targetBits := big.NewInt(1)
	// remain digits for only leading difficulty bits
	targetBits.Lsh(targetBits, targetBitShift)

	randomGenerator := random.NewGenerator()

	return &Hashcash{
		hashFunc:        cfg.HashFunc,
		targetBits:      targetBits,
		randomGenerator: randomGenerator,
	}
}

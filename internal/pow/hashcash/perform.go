package hashcash

import (
	"encoding/binary"
	"log"
	"math"
)

func (h *Hashcash) PerformChallenge(challenge Challenge) UInt64 {
	var counter UInt64
	hasher := h.hashFunc()

	for i := uint64(0); i < math.MaxUint64; i++ {
		attempt := counter[:]
		binary.BigEndian.PutUint64(attempt, i)

		hasher.Reset()

		if checkSolution(challenge.Hash[:], attempt, challenge.Target, hasher) {
			log.Printf("challenge done in %d iterations\n", i)
			break
		}
	}

	return counter
}

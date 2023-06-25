package hashcash

import (
	"encoding/binary"
	"time"
)

func (h *Hashcash) GenerateChallenge() Challenge {
	var nonce UInt64
	binary.BigEndian.PutUint64(nonce[:], h.randomGenerator.UInt64())

	var timestampBytes UInt64
	binary.BigEndian.PutUint64(timestampBytes[:], uint64(time.Now().UnixNano()))

	hasher := h.hashFunc()
	hasher.Write(timestampBytes[:])
	hasher.Write(nonce[:])

	return Challenge{
		Timestamp: timestampBytes,
		Nonce:     nonce,
		Target:    h.targetBits,
		Hash:      hasher.Sum(nil),
	}
}

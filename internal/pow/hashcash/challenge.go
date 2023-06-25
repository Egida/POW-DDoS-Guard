package hashcash

import (
	"math/big"
)

// UInt64 representation in binary.BigEndian notation by a slice of bytes
type UInt64 [8]byte

type Challenge struct {
	Timestamp UInt64
	Nonce     UInt64
	Target    *big.Int
	Hash      []byte
}

package hashcash

import (
	"hash"
	"math/big"
)

func (h *Hashcash) VerifyChallengeResult(challengeResult []byte, nonce, timestamp UInt64) bool {
	hasher := h.hashFunc()

	hasher.Write(timestamp[:])
	hasher.Write(nonce[:])
	challenge := hasher.Sum(nil)

	hasher.Reset()
	return checkSolution(challenge, challengeResult, h.targetBits, hasher)
}

func checkSolution(challenge, challengeRes []byte, target *big.Int, hasher hash.Hash) bool {
	hasher.Write(challenge)
	hasher.Write(challengeRes)

	var result big.Int
	result.SetBytes(hasher.Sum(nil))

	return result.Cmp(target) != 1
}

package bytes

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
)

type ChallengePayload struct {
	Target *big.Int
	Puzzle []byte
}

func (p *ChallengePayload) Serialize() []byte {
	return bytes.Join(
		[][]byte{
			SerializeSlice(p.Target.Bytes()),
			SerializeSlice(p.Puzzle),
		}, nil,
	)
}

func (p *ChallengePayload) Deserialize(r io.Reader) error {
	sliceBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("decerialize target slice: %w", err)
	}

	target := big.NewInt(0)
	target.SetBytes(sliceBytes)
	p.Target = target

	puzzleBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("read puzzle slice: %w", err)
	}
	p.Puzzle = puzzleBytes

	return nil
}

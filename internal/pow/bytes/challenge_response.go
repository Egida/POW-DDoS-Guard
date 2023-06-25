package bytes

import (
	"fmt"
	"io"
)

type ChallengeResponse struct {
	Counter []byte
}

func (p *ChallengeResponse) Serialize() []byte {
	return SerializeSlice(p.Counter)
}

func (p *ChallengeResponse) Deserialize(r io.Reader) error {
	counterBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("deserialize challenge response: %w", err)
	}
	p.Counter = counterBytes

	return nil
}

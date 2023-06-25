package handler

import (
	"crypto/sha256"
	"fmt"
	"net"
	"time"

	"github.com/Goganad/pow_ddos_guard/internal/pow/bytes"
	"github.com/Goganad/pow_ddos_guard/internal/pow/hashcash"
	"github.com/Goganad/pow_ddos_guard/internal/tcp"
)

const (
	readDeadline = 2 * time.Second
)

func WithDDoSGuard(h tcp.Handler) tcp.Handler {
	hashcashPOW := hashcash.New(hashcash.Config{
		Difficulty: 20,
		HashFunc:   sha256.New,
	})

	return func(conn net.Conn) error {
		challenge := hashcashPOW.GenerateChallenge()

		if err := sendChallenge(conn, challenge); err != nil {
			return err
		}

		if err := conn.SetReadDeadline(time.Now().Add(readDeadline)); err != nil {
			return fmt.Errorf("set challenge result read deadline: %v", err)
		}

		if err := verifyChallenge(conn, hashcashPOW, challenge); err != nil {
			return err
		}

		if err := conn.SetDeadline(time.Time{}); err != nil {
			return fmt.Errorf("cant remove previously set deadlines: %v", err)
		}

		return h(conn)
	}
}

func sendChallenge(conn net.Conn, challenge hashcash.Challenge) error {
	challengePayload := bytes.ChallengePayload{
		Target: challenge.Target,
		Puzzle: challenge.Hash,
	}

	if err := bytes.WriteSlice(conn, challengePayload.Serialize()); err != nil {
		return fmt.Errorf("cant write hashcashVerifier client puzzle to tcp conn: %w", err)
	}

	return nil
}

func verifyChallenge(conn net.Conn, hashcash *hashcash.Hashcash, challenge hashcash.Challenge) error {
	var clientResponse bytes.ChallengeResponse
	if err := clientResponse.Deserialize(conn); err != nil {
		return fmt.Errorf("cant read hashcashVerifier client response: %w", err)
	}

	if !hashcash.VerifyChallengeResult(clientResponse.Counter, challenge.Nonce, challenge.Timestamp) {
		return fmt.Errorf("client didnt solve the puzzle")
	}

	return nil
}

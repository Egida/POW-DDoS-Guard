package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net"

	"github.com/Goganad/pow_ddos_guard/internal/pow/bytes"
	"github.com/Goganad/pow_ddos_guard/internal/pow/hashcash"
)

// Default server address is localhost:9991
// Using docker this parameter value will include address and port for server container
var addr = flag.String("addr", ":9991", "")

func main() {
	flag.Parse()

	log.SetFlags(log.Lmicroseconds)

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("net.Dial: %v", err)
	}

	defer func() { _ = conn.Close() }()

	solveChallenge(conn)

	quoteBytes, err := bytes.DeserializeSlice(conn)
	if err != nil {
		log.Fatalf("read quote: %v", err)
	}

	log.Printf("received quote: %q", quoteBytes)
}

func solveChallenge(conn net.Conn) {
	hashCash := hashcash.New(
		hashcash.Config{
			HashFunc: sha256.New,
		})

	var challenge bytes.ChallengePayload

	if err := challenge.Deserialize(conn); err != nil {
		log.Fatalf("read hashcash challenge: %v", err)
	}

	counter := hashCash.PerformChallenge(hashcash.Challenge{
		Target: challenge.Target,
		Hash:   challenge.Puzzle,
	})

	response := bytes.ChallengeResponse{Counter: counter[:]}
	if err := bytes.WriteSlice(conn, response.Serialize()); err != nil {
		log.Fatalf("send hashcash challenge response: %v", err)
	}
}

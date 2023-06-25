package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/Goganad/pow_ddos_guard/assets"
	pow_bytes "github.com/Goganad/pow_ddos_guard/internal/pow/bytes"
	"github.com/Goganad/pow_ddos_guard/internal/tcp"
)

func ReadQuotes(source io.Reader) ([][]byte, error) {
	var quotes [][]byte
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan quotes: %w", err)
	}

	return quotes, nil
}

func RandomQuote() tcp.Handler {
	quotes, err := ReadQuotes(bytes.NewReader(assets.QuotesFile))
	if err != nil {
		log.Fatalf("read quotes: %v", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(conn net.Conn) error {
		randomIndex := r.Intn(len(quotes))

		if err := pow_bytes.WriteSlice(conn, pow_bytes.SerializeSlice(quotes[randomIndex])); err != nil {
			return fmt.Errorf("send quote to client: %w", err)
		}

		return nil
	}
}

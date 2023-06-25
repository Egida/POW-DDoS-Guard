package tcp

import (
	"log"
	"net"
)

type Handler func(conn net.Conn) error

func WithConnectionClosure(h Handler) Handler {
	return func(conn net.Conn) error {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("closing connection error: %v\n", err)
			}
		}()

		return h(conn)
	}
}

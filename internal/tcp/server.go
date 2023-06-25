package tcp

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goganad/pow_ddos_guard/internal/async"
)

type ServerConfig struct {
	Addr          string
	WorkerPoolCfg async.WorkerPoolConfig
}

type Server struct {
	cfg ServerConfig

	workerPool *async.WorkerPool
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		cfg:        cfg,
		workerPool: async.NewWorkerPool(cfg.WorkerPoolCfg),
	}
}

func (s *Server) Serve(ctx context.Context, handler Handler) {
	ctx, cancel := context.WithCancel(ctx)
	listenConfig := net.ListenConfig{}

	l, err := listenConfig.Listen(ctx, "tcp", s.cfg.Addr)
	if err != nil {
		log.Fatalf("can't start serving on port %s: %v", s.cfg.Addr, err)
	}

	log.Printf("server listening on port %s", s.cfg.Addr)

	go s.handleSignals(cancel, l)

	handler = WithConnectionClosure(handler)

	// Accept incoming connections
	for {
		conn, err := l.Accept()
		if err != nil {
			// Check if the error is not due to closing the listener
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("listener.Accept error: %v", err)
			}
			return
		}

		// Handle the connection by worker pool
		s.workerPool.Enqueue(func() {
			if err := handler(conn); err != nil {
				log.Printf("handle tcp connection error: %v", err)
			}
		})
	}
}

// handleSignals listens for SIGTERM signal and performs graceful shutdown
func (s *Server) handleSignals(cancelFunc context.CancelFunc, listener net.Listener) {
	terminationCh := make(chan os.Signal, 1)
	signal.Notify(terminationCh, os.Interrupt, syscall.SIGTERM)

	sig := <-terminationCh
	log.Printf("received signal: %s", sig.String())

	cancelFunc()

	// Close the connection
	if err := listener.Close(); err != nil {
		log.Printf("closing server error: %v", err)
	}

	s.workerPool.Close()
}

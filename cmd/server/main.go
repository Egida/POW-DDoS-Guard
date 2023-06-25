package main

import (
	"context"
	"flag"
	"log"

	"github.com/Goganad/pow_ddos_guard/internal/async"
	"github.com/Goganad/pow_ddos_guard/internal/handler"
	"github.com/Goganad/pow_ddos_guard/internal/tcp"
)

var addr = flag.String("addr", ":9991", "")

func main() {
	log.SetFlags(log.Lmicroseconds)

	flag.Parse()

	srv := tcp.NewServer(tcp.ServerConfig{
		Addr:          *addr,
		WorkerPoolCfg: async.WorkerPoolDefaultConfig(),
	})

	ctx := context.Background()

	srv.Serve(ctx, handler.WithDDoSGuard(
		handler.RandomQuote(),
	))
}

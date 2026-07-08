package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"monolith/cmd/server"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// debug := flag.Bool("debug", false, "sets log level to debug")
	// flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// if *debug {
	// 	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// }

	srv, err := server.NewServer(
		server.WithDB(),
		server.WithDriverDomain(),
		server.WithGeolocationDomain(),
		server.WithAuth(),
	)
	if err != nil {
		log.Printf("Failed to setup server: %v", err)
	}
	err = srv.Start()
	if err != nil {
		log.Fatal("Server failed", err)
	}

	// TODO: graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println("failed to shutdown server")
	}
	log.Println("Shutting down server...")
}

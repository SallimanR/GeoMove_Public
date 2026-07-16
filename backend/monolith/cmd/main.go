package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"monolith/cmd/server"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()

	srv, err := server.NewServer(
		server.WithLogger(logger),
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
		log.Fatal("Server failed: ", err)
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

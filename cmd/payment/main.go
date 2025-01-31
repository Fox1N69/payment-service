package main

import (
	"log"

	"github.com/Fox1N69/iq-testtask/internal/server"
)

func main() {
	srv, err := server.InitializeServer()
	if err != nil {
		log.Fatalf("could not initialize the server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}

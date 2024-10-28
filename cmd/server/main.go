package main

import (
	"log"

	"github.com/teamcutter/atm/internal/server"
)

func main() {
	srv := server.New(server.Config{})

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

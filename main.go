package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/teamcutter/atm/internal/server"
)

func main() {
	srv := server.New(server.Config{})

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.AcceptAndHandle(); err != nil {
			log.Printf("Error in AcceptAndHandle: %v", err)
			srv.Stop()
		}
	}()

	<-sigChan

	log.Println("Stopping server...")
	if err := srv.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
	log.Println("Server stopped gracefully")
}

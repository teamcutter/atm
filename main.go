package main

import (
	"flag"
	"log"

	"github.com/teamcutter/atm/server"
)

func main() {
	port := flag.String("p", "8001", "Port for the server to listen on")
	flag.Parse()

	cfg := server.Config{
		ListenAddr: ":" + *port,
	}

	s := server.New(cfg)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
package main

import (
	"log"

	"github.com/teamcutter/atm/internal/client"
)

func main() {
	serverAddr := "localhost:8001"

	c, err := client.New(serverAddr)
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}
	defer c.Close()

	key := "username"
	value := "john_doe"
	err = c.Set(key, value)
	if err != nil {
		log.Printf("error sending value: %v", err)
	} else {
		log.Printf("Sent %s = %s\n", key, value)
	}

	retrievedValue, err := c.Get(key)
	if err != nil {
		log.Printf("error getting value: %v", err)
	} else {
		log.Printf("Retrieved %s = %s\n", key, retrievedValue)
	}
}

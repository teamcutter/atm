package main

import (
	"fmt"
	"log"

	"github.com/teamcutter/atm/internal/client"
)

func main() {
	// Create a new client and connect to the server
	c, err := client.New("localhost:8001")
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer c.Close()
	
	err = c.Set("greeting", "hello")
	if err != nil {
		log.Printf("error setting value: %v", err)
	} else {
		fmt.Println("Value set successfully")
	}

	val, err := c.Get("greeting")
	if err != nil {
		log.Printf("error getting value: %v", err)
	} else {
		fmt.Printf("Received value: %s\n", val)
	}
}

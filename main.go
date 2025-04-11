package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/teamcutter/atm/server"
)

func main() {
    // Define flags
    pass := flag.String("pass", "", "Password for authentication")
    login := flag.String("login", "", "Login for authentication")
    port := flag.String("p", ":8001", "Port to listen on")

    // Parse flags
    flag.Parse()

    // Validate flags
    if *pass == "" || *login == "" {
        fmt.Println("Error: -pass and -login are required")
        flag.Usage()
        os.Exit(1)
    }

    // Start the server
    srv := server.New(*pass, *login, *port)
    if err := srv.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
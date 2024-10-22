package main

import (
	"log"

	"github.com/teamcutter/atm/internal/server"
)

func main() {
	s := server.New(server.Config{})
	err := s.Start()
	if err != nil {
		panic(err)
	}

	err = s.AcceptAndHandle()
	if err != nil {
		log.Println(err.Error())
	}
	
}

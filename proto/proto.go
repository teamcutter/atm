package proto

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/teamcutter/atm/peers"
)

type Command interface {
	Execute(peer *peers.Peer) error
	String() string
}

type CommandSET struct {
	Key   string
	Value string
}

func (c *CommandSET) Execute(peer *peers.Peer) error {
	log.Printf("Executing SET command: key=%s, value=%s", c.Key, c.Value)

	peer.Set(c.Key, c.Value)
	response := fmt.Sprintf("SET OK: %s = %s\n", c.Key, c.Value)
	return peer.Send(response)
}

func (c *CommandSET) String() string {
	return "SET"
}

type CommandGET struct {
	Key string
}

func (c *CommandGET) Execute(peer *peers.Peer) error {
	log.Printf("Executing GET command: key=%s", c.Key)

	val, err := peer.Get(c.Key)
	if err != nil {
		return err
	}
	response := fmt.Sprintf("VALUE: %s\n", val)
	return peer.Send(response)
}

func (c *CommandGET) String() string {
	return "GET"
}

func parseCommand(msg string) (Command, error) {
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		return nil, errors.New("invalid command format")
	}

	cmdType := strings.ToUpper(parts[0])
	key := parts[1]

	switch cmdType {
	case "SET":
		if len(parts) < 3 {
			return nil, errors.New("SET command requires a key and a value")
		}
		value := parts[2]
		return &CommandSET{Key: key, Value: value}, nil

	case "GET":
		return &CommandGET{Key: key}, nil

	default:
		return nil, errors.New("unknown command")
	}
}

func HandleCommand(msg string, peer *peers.Peer) error {
	cmd, err := parseCommand(msg)
	if err != nil {
		log.Printf("command parsing error: %v", err)
		return err
	}

	err = cmd.Execute(peer)
	if err != nil {
		log.Printf("command execution error: %v", err)
		return err
	}

	log.Printf("Successfully executed command: %s", cmd.String())
	return nil
}

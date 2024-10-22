package proto

import (
	"errors"
	"log"
	"strings"

	"github.com/teamcutter/atm/internal/peers"
)

type Command interface {
	String() string
}

type CommandSET struct {
	Key   string
	Value string
}

func (c *CommandSET) String() string {
	return "SET"
}

type CommandGET struct {
	Key string
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

func HandleCommand(msg string, p *peers.Peer) error {
	cmd, err := parseCommand(string(msg))
	if err != nil {
		return err
	}
	switch c := cmd.(type) {
	case *CommandSET:
		log.Printf("SET key: %v, value: %v", c.Key, c.Value)
		p.Set(c.Key, c.Value)
	case *CommandGET:
		val, err := p.Get(c.Key)
		if err != nil {
			return err
		}
		log.Println(val)
	default:
		log.Printf("Unknown command type: %T", cmd)
	}

	return nil
}

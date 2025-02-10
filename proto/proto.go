package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/teamcutter/atm/peers"
)

// Command constants representing different command types.
const (
	CommandSet = "SET"
	CommandGet = "GET"
	CommandDel = "DEL"
)

// Command defines the interface for all commands.
type Command interface {
	// Execute runs the command on the given peer and returns the response.
	Execute(peer *peers.Peer) (string, error)
	// Serialize converts the command into a byte slice.
	Serialize() ([]byte, error)
	// Deserialize populates the command fields from a byte slice.
	Deserialize(data []byte) error
	// String returns the command type as a string.
	String() string
}

// CommandSET represents the SET command, which stores a key-value pair.
type CommandSET struct {
	Key   string
	Value string
}

// Execute executes the SET command on the given peer.
func (c *CommandSET) Execute(peer *peers.Peer) (string, error) {
	peer.Set(c.Key, c.Value)
	return fmt.Sprintf("SET %s = %s", c.Key, c.Value), nil
}

// Serialize converts the CommandSET into a byte slice.
func (c *CommandSET) Serialize() ([]byte, error) {
	header := []byte(CommandSet) // Command header
	if len(header) != 4 {
		return nil, errors.New("invalid header length")
	}

	// Convert key and value lengths to big-endian byte format
	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	valueLen := uint32(len(c.Value))
	valueLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valueLenBytes, valueLen)

	// Construct the final byte slice
	data := bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key), valueLenBytes, []byte(c.Value)}, nil)
	return data, nil
}

// Deserialize extracts the fields from a byte slice into CommandSET.
func (c *CommandSET) Deserialize(data []byte) error {
	if len(data) < 12 {
		return errors.New("invalid data length")
	}

	header := string(data[:4])
	if header != CommandSet {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandSet, header)
	}

	// Extract key length and value length
	keyLen := binary.BigEndian.Uint32(data[4:8])
	c.Key = string(data[8 : 8+keyLen])

	valueLen := binary.BigEndian.Uint32(data[8+keyLen : 12+keyLen])
	c.Value = string(data[12+keyLen : 12+keyLen+valueLen])

	return nil
}

// String returns the string representation of the SET command.
func (c *CommandSET) String() string {
	return CommandSet
}

// CommandGET represents the GET command, which retrieves a value by key.
type CommandGET struct {
	Key string
}

// Execute executes the GET command and retrieves the value from the peer.
func (c *CommandGET) Execute(peer *peers.Peer) (string, error) {
	val, err := peer.Get(c.Key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("GET %s = %s", c.Key, val), nil
}

// Serialize converts CommandGET into a byte slice.
func (c *CommandGET) Serialize() ([]byte, error) {
	header := []byte(CommandGet)
	if len(header) != 4 {
		return nil, errors.New("invalid header length")
	}

	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	return bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key)}, nil), nil
}

// Deserialize extracts the key from a byte slice into CommandGET.
func (c *CommandGET) Deserialize(data []byte) error {
	if len(data) < 8 {
		return errors.New("invalid data length")
	}

	header := string(data[:4])
	if header != CommandGet {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandGet, header)
	}

	keyLen := binary.BigEndian.Uint32(data[4:8])
	c.Key = string(data[8 : 8+keyLen])

	return nil
}

// String returns the string representation of the GET command.
func (c *CommandGET) String() string {
	return CommandGet
}

// CommandDEL represents the DEL command, which deletes a key from the peer.
type CommandDEL struct {
	Key string
}

// Execute executes the DEL command on the given peer.
func (c *CommandDEL) Execute(peer *peers.Peer) (string, error) {
	val, err := peer.Delete(c.Key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("DEL %s = %s", c.Key, val), nil
}

// Serialize converts CommandDEL into a byte slice.
func (c *CommandDEL) Serialize() ([]byte, error) {
	header := []byte(CommandDel)
	if len(header) != 4 {
		return nil, errors.New("invalid header length")
	}

	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	return bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key)}, nil), nil
}

// Deserialize extracts the key from a byte slice into CommandDEL.
func (c *CommandDEL) Deserialize(data []byte) error {
	if len(data) < 8 {
		return errors.New("invalid data length")
	}

	header := string(data[:4])
	if header != CommandDel {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandDel, header)
	}

	keyLen := binary.BigEndian.Uint32(data[4:8])
	c.Key = string(data[8 : 8+keyLen])

	return nil
}

// String returns the string representation of the DEL command.
func (c *CommandDEL) String() string {
	return CommandDel
}

// HandleCommand processes an incoming command string.
func HandleCommand(msg string, peer *peers.Peer) (string, error) {
	slog.Info("Handling command", "message", msg)

	cmd, err := parseCommand(msg)
	if err != nil {
		slog.Error("Command parsing failed", "error", err)
		return "", fmt.Errorf("command parsing failed: %w", err)
	}

	response, err := cmd.Execute(peer)
	if err != nil {
		slog.Error("Command execution failed", "error", err)
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	slog.Info("Command executed successfully", "response", response)
	return response, nil
}

// parseCommand parses a raw command string into a Command instance.
func parseCommand(msg string) (Command, error) {
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		return nil, errors.New("invalid command format: expected at least 2 parts")
	}

	switch strings.ToUpper(parts[0]) {
	case CommandSet:
		if len(parts) < 3 {
			return nil, fmt.Errorf("SET command requires a key and a value")
		}
		return &CommandSET{Key: parts[1], Value: parts[2]}, nil
	case CommandGet:
		return &CommandGET{Key: parts[1]}, nil
	case CommandDel:
		return &CommandDEL{Key: parts[1]}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", parts[0])
	}
}

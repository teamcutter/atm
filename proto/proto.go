package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Command constants representing different command types.
const (
	CommandSet = "SET"
	CommandGet = "GET"
	CommandDel = "DEL"
)

// Storage defines the interface for key-value storage.
type Storage interface {
	Set(key, value string)
	Get(key string) (string, error)
	Delete(key string) (string, error)
}

// Command defines the interface for all commands.
type Command interface {
	Execute(storage Storage) (string, error)
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	String() string
}

// CommandSET represents the SET command, which stores a key-value pair.
type CommandSET struct {
	Key   string
	Value string
}

func (c *CommandSET) Execute(storage Storage) (string, error) {
	storage.Set(c.Key, c.Value)
	return fmt.Sprintf("SET %s = %s", c.Key, c.Value), nil
}

func (c *CommandSET) Serialize() ([]byte, error) {
	header := []byte(CommandSet)
	if len(header) != 3 {
		return nil, errors.New("invalid header length")
	}

	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	valueLen := uint32(len(c.Value))
	valueLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valueLenBytes, valueLen)

	data := bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key), valueLenBytes, []byte(c.Value)}, nil)
	return data, nil
}

func (c *CommandSET) Deserialize(data []byte) error {
	if len(data) < 11 { // 3 (header) + 4 (keyLen) + 4 (valueLen)
		return errors.New("invalid data length")
	}

	header := string(data[:3])
	if header != CommandSet {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandSet, header)
	}

	keyLen := binary.BigEndian.Uint32(data[3:7]) // Key length: bytes 3-6
	if len(data) < 7+int(keyLen)+4 { // Check if enough data for key + valueLen
		return errors.New("insufficient data for key and value length")
	}
	c.Key = string(data[7 : 7+keyLen])

	valueLen := binary.BigEndian.Uint32(data[7+keyLen : 11+keyLen]) // Value length: bytes after key
	if len(data) < 11+int(keyLen)+int(valueLen) { // Check if enough data for value
		return errors.New("insufficient data for value")
	}
	c.Value = string(data[11+keyLen : 11+keyLen+valueLen])

	return nil
}

func (c *CommandSET) String() string {
	return CommandSet
}

// CommandGET represents the GET command, which retrieves a value by key.
type CommandGET struct {
	Key string
}

func (c *CommandGET) Execute(storage Storage) (string, error) {
	val, err := storage.Get(c.Key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("GET %s = %s", c.Key, val), nil
}

func (c *CommandGET) Serialize() ([]byte, error) {
	header := []byte(CommandGet)
	if len(header) != 3 {
		return nil, errors.New("invalid header length")
	}

	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	return bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key)}, nil), nil
}

func (c *CommandGET) Deserialize(data []byte) error {
	if len(data) < 7 { // 3 (header) + 4 (keyLen)
		return errors.New("invalid data length")
	}

	header := string(data[:3])
	if header != CommandGet {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandGet, header)
	}

	keyLen := binary.BigEndian.Uint32(data[3:7])
	if len(data) < 7+int(keyLen) {
		return errors.New("insufficient data for key")
	}
	c.Key = string(data[7 : 7+keyLen])

	return nil
}

func (c *CommandGET) String() string {
	return CommandGet
}

// CommandDEL represents the DEL command, which deletes a key from the storage.
type CommandDEL struct {
	Key string
}

func (c *CommandDEL) Execute(storage Storage) (string, error) {
	val, err := storage.Delete(c.Key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("DEL %s = %s", c.Key, val), nil
}

func (c *CommandDEL) Serialize() ([]byte, error) {
	header := []byte(CommandDel)
	if len(header) != 3 {
		return nil, errors.New("invalid header length")
	}

	keyLen := uint32(len(c.Key))
	keyLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyLenBytes, keyLen)

	return bytes.Join([][]byte{header, keyLenBytes, []byte(c.Key)}, nil), nil
}

func (c *CommandDEL) Deserialize(data []byte) error {
	if len(data) < 7 { // 3 (header) + 4 (keyLen)
		return errors.New("invalid data length")
	}

	header := string(data[:3])
	if header != CommandDel {
		return fmt.Errorf("invalid header: expected %s, got %s", CommandDel, header)
	}

	keyLen := binary.BigEndian.Uint32(data[3:7])
	if len(data) < 7+int(keyLen) {
		return errors.New("insufficient data for key")
	}
	c.Key = string(data[7 : 7+keyLen])

	return nil
}

func (c *CommandDEL) String() string {
	return CommandDel
}
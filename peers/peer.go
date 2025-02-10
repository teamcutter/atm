package peers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Message represents a message sent between peers.
type Message struct {
	Sender *Peer // The sender of the message.
	Cmd    []byte // The command or message content.
}

// Peer represents a connected client in the network.
type Peer struct {
	storage sync.Map   // Thread-safe key-value storage for peer-specific data.
	conn    net.Conn   // Network connection associated with the peer.
	msgChan chan Message // Channel for sending messages to the server.
}

// New creates and returns a new Peer instance.
func New(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

// Receive listens for incoming messages from the peer and sends them to msgChan.
// It continuously reads data until an error occurs (e.g., connection closed).
func (p *Peer) Receive() error {
	reader := bufio.NewReader(p.conn)
	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			p.conn.Close()
			return fmt.Errorf("failed to read cmd: %w", err)
		}
		p.msgChan <- Message{
			Sender: p,
			Cmd:    []byte(strings.TrimSpace(cmd)),
		}
	}
}

// Send writes a message to the peer's connection.
func (p *Peer) Send(msg string) error {
	_, err := p.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message to peer: %w", err)
	}
	return nil
}

// Close terminates the peer's connection.
func (p *Peer) Close() {
	p.conn.Close()
}

// Set stores a key-value pair in the peer's storage.
func (p *Peer) Set(key string, value string) {
	p.storage.Store(key, value)
}

// Get retrieves a value from the peer's storage by key.
// Returns an error if the key does not exist.
func (p *Peer) Get(key string) (string, error) {
	val, ok := p.storage.Load(key)
	if !ok {
		return "", errors.New("no record with such key")
	}
	return val.(string), nil
}

// Delete removes a key from the peer's storage and returns its value.
// Returns an error if the key does not exist.
func (p *Peer) Delete(key string) (string, error) {
	val, ok := p.storage.Load(key)
	if !ok {
		return "", errors.New("no record with such key")
	}
	p.storage.Delete(key)
	return val.(string), nil
}

// Clear removes all stored key-value pairs from the peer's storage.
func (p *Peer) Clear() {
	p.storage.Clear()
}

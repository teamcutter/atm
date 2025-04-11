package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Client represents a client that communicates with the ATM server.
type Client struct {
	conn       net.Conn
	addr       string
	login      string
	password   string
	responses  chan string
	listenerWG sync.WaitGroup
	closed     bool
	mu         sync.Mutex
}

// New creates a new Client instance and connects to the server.
func New(addr, login, password string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := &Client{
		conn:      conn,
		addr:      addr,
		login:     login,
		password:  password,
		responses: make(chan string, 10),
	}

	// Authenticate
	if err := client.authenticate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Start listening for responses
	client.listenerWG.Add(1)
	go client.listen()

	return client, nil
}

// authenticate sends login:password and checks the server's response.
func (c *Client) authenticate() error {
	auth := fmt.Sprintf("%s:%s\n", c.login, c.password)
	if _, err := c.conn.Write([]byte(auth)); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	reader := bufio.NewReader(c.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}
	response = strings.TrimSpace(response)
	if response != "OK" {
		return fmt.Errorf("invalid auth response: %s", response)
	}
	return nil
}

// listen continuously reads server responses and sends them to the responses channel.
func (c *Client) listen() {
	defer c.listenerWG.Done()
	reader := bufio.NewReader(c.conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			c.mu.Lock()
			if !c.closed {
				log.Printf("Failed to read response: %v", err)
				c.responses <- fmt.Sprintf("ERROR: %v", err)
				c.Close()
			}
			c.mu.Unlock()
			return
		}
		c.responses <- strings.TrimSpace(response)
	}
}

// Set sends a SET command to store a key-value pair.
func (c *Client) Set(key, value string) (string, error) {
	var cmd bytes.Buffer
	cmd.WriteString("SET") // 3-byte header
	binary.Write(&cmd, binary.BigEndian, uint32(len(key)))
	cmd.WriteString(key)
	binary.Write(&cmd, binary.BigEndian, uint32(len(value)))
	cmd.WriteString(value)
	cmd.WriteByte('\n')

	return c.sendCommand(cmd.Bytes())
}

// Get sends a GET command to retrieve a value by key.
func (c *Client) Get(key string) (string, error) {
	var cmd bytes.Buffer
	cmd.WriteString("GET") // 3-byte header
	binary.Write(&cmd, binary.BigEndian, uint32(len(key)))
	cmd.WriteString(key)
	cmd.WriteByte('\n')

	response, err := c.sendCommand(cmd.Bytes())
	if err != nil {
		return "", err
	}
	// Parse the response to extract the value
	parts := strings.Split(response, "=")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid GET response format: %s", response)
	}
	return strings.TrimSpace(parts[1]), nil
}

// Del sends a DEL command to delete a key.
func (c *Client) Del(key string) (string, error) {
	var cmd bytes.Buffer
	cmd.WriteString("DEL") // 3-byte header
	binary.Write(&cmd, binary.BigEndian, uint32(len(key)))
	cmd.WriteString(key)
	cmd.WriteByte('\n')

	response, err := c.sendCommand(cmd.Bytes())
	if err != nil {
		return "", err
	}
	// Parse the response to extract the deleted value
	parts := strings.Split(response, "=")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid DEL response format: %s", response)
	}
	return strings.TrimSpace(parts[1]), nil
}

// sendCommand sends a command and waits for the response.
func (c *Client) sendCommand(data []byte) (string, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return "", errors.New("client is closed")
	}
	_, err := c.conn.Write(data)
	c.mu.Unlock()
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	// Wait for response
	response := <-c.responses
	if strings.HasPrefix(response, "ERROR:") {
		return "", fmt.Errorf("server error: %s", response)
	}
	return response, nil
}

// Close terminates the client connection.
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	err := c.conn.Close()
	close(c.responses)
	c.mu.Unlock()
	c.listenerWG.Wait()
	return err
}
package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	mu   sync.Mutex
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) send(msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (c *Client) receive() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	reader := bufio.NewReader(c.conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}
	return strings.TrimSpace(message), nil
}

func (c *Client) Get(key string) (string, error) {
	if err := c.send(fmt.Sprintf("GET %s", key)); err != nil {
		return "", fmt.Errorf("error sending GET request: %v", err)
	}

	resp, err := c.receive()
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if strings.HasPrefix(resp, "VALUE: ") {
		return strings.TrimPrefix(resp, "VALUE: "), nil
	}
	return "", fmt.Errorf("unexpected response: %s", resp)
}

func (c *Client) Set(key string, value string) error {
	if err := c.send(fmt.Sprintf("SET %s %s", key, value)); err != nil {
		return fmt.Errorf("error sending SET request: %v", err)
	}

	resp, err := c.receive()
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	if strings.HasPrefix(resp, "SET OK:") {
		return nil
	}
	return fmt.Errorf("unexpected response: %s", resp)
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

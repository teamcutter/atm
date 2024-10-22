package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) send(msg string) error {
	_, err := c.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (c *Client) receive() (string, error) {
	reader := bufio.NewReader(c.conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}
	return strings.TrimSpace(message), nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (c *Client) Get(key string) (string, error) {
	err := c.send(fmt.Sprintf("GET %s\n", key))
	if err != nil {
		return "", fmt.Errorf("error sending GET request: %v", err)
	}

	resp, err := c.receive()
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if strings.HasPrefix(resp, "VALUE: ") {
		return strings.TrimPrefix(resp, "VALUE: "), nil
	}
	return "", fmt.Errorf("unexpected resp: %s", resp)
}

func (c *Client) Set(key string, value string) error {
	err := c.send(fmt.Sprintf("SET %s %s\n", key, value))
	if err != nil {
		return fmt.Errorf("error sending SET request: %v", err)
	}

	resp, err := c.receive()
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	if strings.HasPrefix(resp, "SET OK:") {
		return nil
	}
	return fmt.Errorf("unexpected resp: %s", resp)
}

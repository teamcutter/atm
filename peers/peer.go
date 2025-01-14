package peers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Message struct {
	Sender *Peer
	Cmd    []byte
}

type Peer struct {
	storage sync.Map
	conn    net.Conn
	msgChan chan Message
}

func New(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

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

func (p *Peer) Send(msg string) error {
	_, err := p.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message to peer: %w", err)
	}
	return nil
}

func (p *Peer) Close() {
	p.conn.Close()
}

func (p *Peer) Set(key string, value string) {
	p.storage.Store(key, value)
}

func (p *Peer) Get(key string) (string, error) {
	val, ok := p.storage.Load(key)
	if !ok {
		return "", errors.New("no record with such key")
	}
	return val.(string), nil
}

func (p *Peer) Delete(key string) (string, error) {
	val, ok := p.storage.Load(key)
	if !ok {
		return "", errors.New("no record with such key")
	}

	p.storage.Delete(key)
	return val.(string), nil
}

func (p *Peer) Clear() {
	p.storage.Clear()
}

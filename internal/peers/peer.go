package peers

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Message struct {
	Sender  *Peer
	Content []byte
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
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			p.conn.Close()
			if err == io.EOF {
				return fmt.Errorf("connection %s closed: eof", p.conn.RemoteAddr())
			}
			return err
		}
		p.msgChan <- Message{
			Sender:  p,
			Content: buf[:n],
		}
	}
}

func (p *Peer) Send(msg string) error {
	_, err := p.conn.Write([]byte(msg + "\n")) // Send message with a newline delimiter
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

func (p *Peer) Clear() {
	p.storage.Clear()
}

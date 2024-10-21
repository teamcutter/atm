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
	data    sync.Map
	conn    net.Conn
	msgChan chan Message
}

func New(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

func (p *Peer) Listen() error {
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

func (p *Peer) Close() {
	p.conn.Close()
}

func (p *Peer) Set(key string, value string) {
	p.data.Store(key, value)
}

func (p *Peer) Get(key string) (string, error) {
	val, ok := p.data.Load(key)
	if !ok {
		return "", errors.New("no record with such key")
	}
	return val.(string), nil
}

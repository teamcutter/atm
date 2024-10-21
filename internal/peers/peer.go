package peers

import (
	"fmt"
	"io"
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan []byte
}

func New(conn net.Conn, msgChan chan []byte) *Peer {
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
		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		p.msgChan <- msgBuf
	}
}

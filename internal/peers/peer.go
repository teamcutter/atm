package peers

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Peer struct {
	conn net.Conn
}

func New(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

func (p *Peer) Listen() error {
	for {
		msg := make([]byte, 1024)
		_, err := p.conn.Read(msg)
		if err != nil {
			p.conn.Close()
			if err == io.EOF {
				return fmt.Errorf("connection %s closed: eof", p.conn.RemoteAddr())
			}
			return err
		}
		
		log.Println(string(msg))
	}
}

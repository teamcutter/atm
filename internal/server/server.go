package server

import (
	"log"
	"net"
	"sync"

	"github.com/teamcutter/atm/internal/peers"
)

const defaultListenAddr = ":8000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers sync.Map
	ln    net.Listener
}

func New(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config: cfg,
		peers:  sync.Map{},
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}

func (s *Server) AcceptAndListen() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := peers.New(conn)
	s.peers.Store(peer, true)
	log.Printf("new connection: %s", conn.RemoteAddr())

	errChan := make(chan error, 1)

	go func() {
		errChan <- peer.Listen()
	}()

	err := <-errChan
	if err != nil {
		log.Println(err.Error())
		s.peers.Delete(peer)
		return
	}
}

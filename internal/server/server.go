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
	peers   sync.Map
	ln      net.Listener
	msgChan chan []byte
	errChan chan error
}

func New(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:  cfg,
		peers:   sync.Map{},
		msgChan: make(chan []byte),
		errChan: make(chan error, 1),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.listen()

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
	defer conn.Close()

	peer := peers.New(conn, s.msgChan)
	s.peers.Store(peer, true)
	log.Printf("new connection: %s", conn.RemoteAddr())

	if err := peer.Listen(); err != nil {
		s.errChan <- err
		s.peers.Delete(peer)
		return
	}
}

func (s *Server) listen() {
	for {
		select {
		case msg := <-s.msgChan:
			log.Println(string(msg))
		case err := <-s.errChan:
			log.Println(err.Error())
		}
	}
}

package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/teamcutter/atm/internal/peers"
	"github.com/teamcutter/atm/internal/proto"
)

const defaultListenAddr = ":8001"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers   sync.Map
	ln      net.Listener
	msgChan chan peers.Message
	errChan chan error
}

func New(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:  cfg,
		peers:   sync.Map{},
		msgChan: make(chan peers.Message),
		errChan: make(chan error, 1),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	log.Println("Starting server...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		if err := s.acceptAndHandle(); err != nil {
			log.Printf("Error in acceptAndHandle: %v", err)
			s.stop()
		}
	}()

	go s.listen()

	<-sigChan
	log.Println("Stopping server...")
	return s.stop()
}

func (s *Server) acceptAndHandle() error {
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

	if err := peer.Receive(); err != nil {
		s.errChan <- err
		s.peers.Delete(peer)
		return
	}
}

func (s *Server) listen() {
	for {
		select {
		case msg := <-s.msgChan:
			err := proto.HandleCommand(string(msg.Cmd), msg.Sender)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		case err := <-s.errChan:
			log.Printf("error: %v", err)
		}
	}
}

func (s *Server) stop() error {
	if s.ln != nil {
		if err := s.ln.Close(); err != nil {
			return err
		}
	}

	s.peers.Range(func(key, value interface{}) bool {
		peer := key.(*peers.Peer)
		peer.Close()
		peer.Clear()
		s.peers.Delete(key)
		return true
	})

	close(s.msgChan)
	close(s.errChan)

	return nil
}

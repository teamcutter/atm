package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/teamcutter/atm/proto"
)

// Server represents a TCP server that handles a single authenticated connection.
type Server struct {
	listenAddr string
	ln         net.Listener
	conn       net.Conn
	login      string
	password   string
	storage    sync.Map
}

func New(password, login, addr string) *Server {
	return &Server{
		listenAddr: addr,
		login:      login,
		password:   password,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.listenAddr, err)
	}
	s.ln = ln
	log.Printf("Server listening on %s", s.listenAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	conn, err := s.ln.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}
	s.conn = conn
	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	if err := s.conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		s.conn.Close()
		return fmt.Errorf("failed to set auth deadline: %w", err)
	}

	if err := s.authenticate(); err != nil {
		log.Printf("Authentication failed: %v", err)
		s.conn.Close()
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := s.conn.SetReadDeadline(time.Time{}); err != nil {
		s.conn.Close()
		return fmt.Errorf("failed to reset deadline: %w", err)
	}

	log.Println("Client authenticated successfully")

	go s.processCommands()

	<-sigChan
	log.Println("Stopping server...")
	return s.Stop()
}

func (s *Server) authenticate() error {
	reader := bufio.NewReader(s.conn)
	auth, err := reader.ReadString('\n')
	if err != nil {
		s.conn.Write([]byte("ERROR: failed to read credentials, please send login:password\n"))
		return fmt.Errorf("failed to read auth: %w", err)
	}
	auth = strings.TrimSpace(auth)
	expected := fmt.Sprintf("%s:%s", s.login, s.password)
	log.Printf("Received auth: %q, expected: %q", auth, expected)
	if auth != expected {
		s.conn.Write([]byte("ERROR: invalid login or password\n"))
		return fmt.Errorf("invalid auth: got %q, expected %q", auth, expected)
	}
	_, err = s.conn.Write([]byte("OK\n"))
	if err != nil {
		return fmt.Errorf("failed to send OK: %w", err)
	}
	return nil
}

func (s *Server) processCommands() {
	reader := bufio.NewReader(s.conn)
	for {
		cmdData, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Failed to read command: %v", err)
			s.Stop()
			return
		}
		cmdData = cmdData[:len(cmdData)-1] // Trim newline
		log.Printf("Received command: %x", cmdData)

		var cmd proto.Command
		switch string(cmdData[:3]) { // Check 3-byte header
		case proto.CommandSet:
			cmd = &proto.CommandSET{}
		case proto.CommandGet:
			cmd = &proto.CommandGET{}
		case proto.CommandDel:
			cmd = &proto.CommandDEL{}
		default:
			log.Printf("Unknown command header: %s", string(cmdData[:3]))
			s.conn.Write([]byte("ERROR: unknown command\n"))
			continue
		}

		if err := cmd.Deserialize(cmdData); err != nil {
			log.Printf("Failed to deserialize command: %v", err)
			s.conn.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
			continue
		}

		response, err := cmd.Execute(s)
		if err != nil {
			log.Printf("Command execution failed: %v", err)
			s.conn.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
			continue
		}

		log.Printf("Sending response: %s", response)
		if _, err := s.conn.Write([]byte(response + "\n")); err != nil {
			log.Printf("Failed to send response: %v", err)
			s.Stop()
			return
		}
	}
}

func (s *Server) Stop() error {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

func (s *Server) Set(key, value string) {
	s.storage.Store(key, value)
}

func (s *Server) Get(key string) (string, error) {
	val, ok := s.storage.Load(key)
	if !ok {
		return "", fmt.Errorf("no record with such key")
	}
	return val.(string), nil
}

func (s *Server) Delete(key string) (string, error) {
	val, ok := s.storage.Load(key)
	if !ok {
		return "", fmt.Errorf("no record with such key")
	}
	s.storage.Delete(key)
	return val.(string), nil
}
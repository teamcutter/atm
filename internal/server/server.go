package server

import "net"

const defaultListenAddr = ":8000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln net.Listener
}

func New(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config: cfg,
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

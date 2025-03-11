package server

import (
	"log/slog"
	"net"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/config"
)

type Server struct {
	lst         net.Listener
	database    Database
	log         *slog.Logger
	bufSize     int
	idleTimeout time.Duration
}

type Database interface {
	Execute(source string) string
}

func New(
	cfg *config.Config,
	database Database,
	log *slog.Logger) (*Server, error) {
	if cfg == nil {
		panic("config is nil")
	}
	if database == nil {
		panic("database is nil")
	}
	if log == nil {
		panic("logger is nil")
	}

	lst, err := net.Listen("tcp", cfg.Network.Address)
	if err != nil {
		return nil, err
	}

	log.Info("started listening", slog.String("address", cfg.Network.Address))
	return &Server{
		lst:         lst,
		database:    database,
		log:         log,
		bufSize:     cfg.Network.MaxMessageSize,
		idleTimeout: cfg.Network.IdleTimeout,
	}, nil
}

func (s *Server) Run() error {
	for {
		conn, err := s.lst.Accept()
		if err != nil {
			return err
		}

		go s.handler(conn)
	}
}

func (s *Server) handler(conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			s.log.Error("panic recovered", slog.Any("error", v))
		}
		conn.Close()
	}()

	var err error
	var n int
	var response string
	buf := make([]byte, s.bufSize)
	for {
		err = conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		if err != nil {
			s.log.Error("set read deadline failure", slog.Any("error", err))
			break
		}

		n, err = conn.Read(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				s.log.Warn("idle connection")
				break
			}
			s.log.Error("read failure", slog.Any("error", err))
			break
		}

		response = s.database.Execute(string(buf[:n]))
		if _, err = conn.Write([]byte(response + "\n")); err != nil {
			s.log.Error("write failure", slog.Any("error", err))
			break
		}
	}
}

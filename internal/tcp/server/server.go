package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/concurrency"
)

type Server struct {
	listener    net.Listener
	bufferSize  int
	idleTimeout time.Duration
	logic       func([]byte) ([]byte, error)
	semaphore   *concurrency.Semaphore
	log         *slog.Logger
}

const (
	defaultBufferSize = 4 << 10
)

//go:generate mockery --name=Database --case=snake --inpackage --inpackage-suffix --with-expecter
type Database interface {
	Execute(source string) string
}

func NewServer(
	address string,
	maxMessageSize int,
	log *slog.Logger, opts ...ServerOption) (*Server, error) {
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	log.Info("started listening", slog.String("address", address))

	server := &Server{
		listener:   listener,
		bufferSize: maxMessageSize,
		log:        log,
	}

	for _, opt := range opts {
		opt(server)
	}

	if server.bufferSize == 0 {
		server.bufferSize = defaultBufferSize
	}

	return server, nil
}

func (s *Server) Run(ctx context.Context, logic func([]byte) ([]byte, error)) error {
	done := make(chan struct{})
	s.logic = logic

	defer func() {
		if err := s.listener.Close(); err != nil {
			s.log.Error("failed to close listener", slog.Any("error", err))
		}
	}()

	go func() {
		defer close(done)
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					s.log.Error("failed to accept connection", slog.Any("error", err))
				}
				return
			}

			go s.recoverer(s.clientsLimiter(s.handler))(ctx, connection)
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Info("stopping server")
		return nil
	case <-done:
		s.log.Error("unexpected server error")
		return errors.New("server stopped accepting connections")
	}
}

func (s *Server) handler(ctx context.Context, connection net.Conn) {
	s.log.Debug("new connection", slog.String("remote", connection.RemoteAddr().String()))

	defer func() {
		if err := connection.Close(); err != nil {
			s.log.Error("failed to close connection", slog.Any("error", err))
		}
	}()

	var err error
	var n int
	buffer := make([]byte, s.bufferSize)
	for {
		if ctx.Err() != nil {
			return
		}

		if s.idleTimeout != 0 {
			if err = connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.log.Error("set read deadline failure", slog.Any("error", err))
				return
			}
		}
		n, err = common.Read(connection, buffer)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				s.log.Warn("idle connection", slog.Any("error", err))
				return
			}
			s.log.Error("read failure", slog.Any("error", err))
			return
		}

		response, err := s.logic(buffer[:n])
		if err != nil {
			s.log.Error("failed to execute logic", slog.Any("error", err))
			return
		}
		if s.idleTimeout != 0 {
			if err = connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.log.Error("set write deadline failure", slog.Any("error", err))
				return
			}
		}
		if _, err = common.Write(connection, response); err != nil {
			s.log.Error("write failure", slog.Any("error", err))
			return
		}
	}
}

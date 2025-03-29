package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/concurrency"
)

type Server struct {
	listener    net.Listener
	bufferSize  int
	idleTimeout time.Duration
	logic       func([]byte) ([]byte, error)
	semaphore   *concurrency.Semaphore
	wg          sync.WaitGroup
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
	s.logic = logic
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			return err
		}

		s.wg.Add(1)
		go s.recoverer(s.clientsLimiter(s.handler))(ctx, connection)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.listener.Close(); err != nil {
		s.log.Error("failed to close listener", slog.Any("error", err))
	}

	doneConnections := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(doneConnections)
	}()

	select {
	case <-ctx.Done():
		s.log.Info("force shutdown")
	case <-doneConnections:
		s.log.Info("all connections closed")
	}
}

func (s *Server) handler(ctx context.Context, conn net.Conn) {
	s.log.Debug("new connection", slog.String("remote", conn.RemoteAddr().String()))

	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Error("failed to close connection", slog.Any("error", err))
		}
		s.wg.Done()
	}()

	var err error
	var n int
	buffer := make([]byte, s.bufferSize)
	for {
		if ctx.Err() != nil {
			return
		}

		if s.idleTimeout != 0 {
			err = conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
			if err != nil {
				s.log.Error("set read deadline failure", slog.Any("error", err))
				break
			}
		}

		n, err = conn.Read(buffer)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				s.log.Warn("idle connection", slog.Any("error", err))
				break
			}
			s.log.Error("read failure", slog.Any("error", err))
			break
		}

		response, err := s.logic(buffer[:n])
		if err != nil {
			s.log.Error("failed to execute logic", slog.Any("error", err))
			break
		}

		if _, err = conn.Write(response); err != nil {
			s.log.Error("write failure", slog.Any("error", err))
			break
		}
	}
}

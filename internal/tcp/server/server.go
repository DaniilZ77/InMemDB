package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/concurrency"
)

type Server struct {
	lst         net.Listener
	database    Database
	log         *slog.Logger
	bufSize     int
	idleTimeout time.Duration
	semaphore   *concurrency.Semaphore
	wg          sync.WaitGroup
}

//go:generate mockery --name=Database --case=snake --inpackage --inpackage-suffix --with-expecter
type Database interface {
	Execute(source string) string
}

func NewServer(
	address string,
	maxMessageSize int,
	idleTimeout time.Duration,
	maxConnections int,
	database Database,
	log *slog.Logger) (*Server, error) {
	if database == nil {
		return nil, errors.New("database is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	lst, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	log.Info("started listening", slog.String("address", address))

	return &Server{
		lst:         lst,
		database:    database,
		log:         log,
		bufSize:     maxMessageSize,
		idleTimeout: idleTimeout,
		semaphore:   concurrency.NewSemaphore(maxConnections),
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	for {
		conn, err := s.lst.Accept()
		if err != nil {
			return err
		}

		s.wg.Add(1)
		go s.recoverer(s.clientsLimiter(s.handler))(ctx, conn)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.lst.Close(); err != nil {
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
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Error("failed to close connection", slog.Any("error", err))
		}
		s.wg.Done()
	}()

	var err error
	var n int
	var resp string
	buf := make([]byte, s.bufSize)
	for {
		if ctx.Err() != nil {
			return
		}

		err = conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		if err != nil {
			s.log.Error("set read deadline failure", slog.Any("error", err))
			break
		}

		n, err = conn.Read(buf)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				s.log.Warn("idle connection", slog.Any("error", err))
				break
			}
			s.log.Error("read failure", slog.Any("error", err))
			break
		}

		resp = s.database.Execute(string(buf[:n]))
		if _, err = fmt.Fprintln(conn, resp); err != nil {
			s.log.Error("write failure", slog.Any("error", err))
			break
		}
	}
}

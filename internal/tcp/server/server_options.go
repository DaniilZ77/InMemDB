package server

import (
	"time"

	"github.com/DaniilZ77/InMemDB/internal/concurrency"
)

type ServerOption func(*Server)

func WithIdleTimeout(idleTimeout time.Duration) ServerOption {
	return func(s *Server) {
		s.idleTimeout = idleTimeout
	}
}

func WithMaxConnections(maxConnections int) ServerOption {
	return func(s *Server) {
		s.semaphore = concurrency.NewSemaphore(maxConnections)
	}
}

func WithMaxMessageSize(maxMessageSize int) ServerOption {
	return func(s *Server) {
		s.bufferSize = maxMessageSize
	}
}

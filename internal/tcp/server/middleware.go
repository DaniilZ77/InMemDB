package server

import (
	"context"
	"log/slog"
	"net"
)

type handler func(context.Context, net.Conn)

func (s *Server) clientsLimiter(next handler) handler {
	return func(ctx context.Context, conn net.Conn) {
		if s.semaphore != nil {
			s.semaphore.Acquire()
			defer s.semaphore.Release()
		}

		next(ctx, conn)
	}
}

func (s *Server) recoverer(next handler) handler {
	return func(ctx context.Context, conn net.Conn) {
		defer func() {
			if v := recover(); v != nil {
				s.log.Error("panic recovered", slog.Any("error", v))
			}
		}()

		next(ctx, conn)
	}
}

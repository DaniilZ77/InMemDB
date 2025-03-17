package server

import (
	"log/slog"
	"net"
)

func (s *Server) clientsLimiter(next func(conn net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		s.semaphore.Acquire()
		defer s.semaphore.Release()

		next(conn)
	}
}

func (s *Server) recoverer(next func(net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		defer func() {
			if v := recover(); v != nil {
				s.log.Error("panic recovered", slog.Any("error", v))
			}
		}()

		next(conn)
	}
}

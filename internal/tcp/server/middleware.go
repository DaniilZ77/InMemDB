package server

import (
	"log/slog"
	"net"
)

func (s *Server) clientsLimiter(next func(conn net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		s.mu.Lock()
		if s.clients >= s.maxClients {
			s.mu.Unlock()
			conn.Close()
			return
		}
		s.clients++
		s.mu.Unlock()

		s.log.Debug("amount of clients increased", slog.Int("clients", s.clients))

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

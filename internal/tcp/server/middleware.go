package server

import (
	"log/slog"
	"net"
)

func (s *Server) clientsLimiter(next func(conn net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		s.condition.L.Lock()
		for s.clients >= s.maxClients {
			s.condition.Wait()
		}
		s.clients++
		s.condition.L.Unlock()

		defer func() {
			s.condition.L.Lock()
			s.clients--
			s.condition.L.Unlock()

			s.condition.Signal()
		}()

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

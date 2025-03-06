package tcp

import (
	"context"
	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/domain/models"
	"log/slog"
	"net"
)

type Server struct {
	lst      net.Listener
	log      *slog.Logger
	parser   Parser
	database Database
	bufSize  int
}

type Parser interface {
	Parse(input string) (*models.Command, error)
}

type Database interface {
	Set(ctx context.Context, cmd models.SetCommand) error
	Get(ctx context.Context, cmd models.GetCommand) ([]byte, error)
	Del(ctx context.Context, cmd models.DelCommand) error
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	parser Parser,
	database Database) (*Server, error) {
	lst, err := net.Listen("tcp", cfg.DbPort)
	if err != nil {
		return nil, err
	}

	log.Info("started listening", slog.String("port", cfg.DbPort))
	return &Server{lst, log, parser, database, cfg.BufSize}, nil
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
	var err error
	var n int
	var command *models.Command
	for {
		buf := make([]byte, s.bufSize)
		n, err = conn.Read(buf)
		if err != nil {
			s.log.Error("read failure", slog.Any("error", err))
			break
		}
		command, err = s.parser.Parse(string(buf[:n]))
		if err != nil {
			s.log.Debug("parse failure", slog.Any("error", err))
			if _, err = conn.Write([]byte("INVALID\n")); err != nil {
				s.log.Error("write failure", slog.Any("error", err))
			}
			continue
		}

		switch command.Type {
		case models.GET:
			s.log.Info("get", slog.Any("command", *command))
			if _, err = conn.Write([]byte("OK\n")); err != nil {
				s.log.Error("write failure", slog.Any("error", err))
			}
		case models.SET:
			s.log.Info("set", slog.Any("command", *command))
			if _, err = conn.Write([]byte("OK\n")); err != nil {
				s.log.Error("write failure", slog.Any("error", err))
			}
		case models.DEL:
			s.log.Info("del", slog.Any("command", *command))
			if _, err = conn.Write([]byte("OK\n")); err != nil {
				s.log.Error("write failure", slog.Any("error", err))
			}
		default:
			s.log.Error("invalid command type", slog.Any("command", *command))
		}
	}
}

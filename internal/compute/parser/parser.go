package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Parser struct {
	log *slog.Logger
}

func NewParser(log *slog.Logger) (*Parser, error) {
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	return &Parser{log: log}, nil
}

func (p *Parser) Parse(source string) (*Command, error) {
	tokens := strings.Fields(source)
	if len(tokens) == 0 {
		p.log.Warn("empty command")
		return nil, fmt.Errorf("%w: empty command", ErrInvalidCommand)
	}

	token := strings.ToLower(tokens[0])
	if commandType, ok := keywords[token]; !ok {
		p.log.Warn("bad command type", slog.String("command type", token))
		return nil, fmt.Errorf("%w: bad command type", ErrInvalidCommand)
	} else {
		return p.parseArgs(commandType, tokens[1:])
	}
}

func (p *Parser) parseArgs(commandType CommandType, tokens []string) (*Command, error) {
	if len(tokens) != commandType.argsCount() {
		p.log.Warn("bad amount of args", slog.Int("args", len(tokens)), slog.Int("expected", commandType.argsCount()))
		return nil, fmt.Errorf("%w: bad amount of args", ErrInvalidCommand)
	}

	return &Command{
		Type: commandType,
		Args: tokens,
	}, nil
}

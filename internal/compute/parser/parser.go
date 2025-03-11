package parser

import (
	"fmt"
	"log/slog"
	"strings"
)

type Parser struct {
	log *slog.Logger
}

func NewParser(log *slog.Logger) *Parser {
	if log == nil {
		panic("logger is nil")
	}

	return &Parser{log: log}
}

func (p *Parser) Parse(source string) (*Command, error) {
	tokens := strings.Fields(source)
	if len(tokens) == 0 {
		p.log.Warn("empty command")
		return nil, fmt.Errorf("%w: empty command", ErrInvalidCommand)
	}

	token := tokens[0]
	if commandType, ok := keywords[token]; !ok {
		p.log.Warn("bad command type", slog.String("command type", token))
		return nil, fmt.Errorf("%w: bad command type", ErrInvalidCommand)
	} else {
		return p.parseArgs(commandType, tokens[1:])
	}
}

func (p *Parser) parseArgs(commandType commandType, tokens []string) (*Command, error) {
	if len(tokens) != commandType.argsCount() {
		p.log.Warn("bad args", slog.Any("args", tokens))
		return nil, fmt.Errorf("%w: bad args", ErrInvalidCommand)
	}

	return &Command{
		Type: commandType,
		Args: tokens,
	}, nil
}

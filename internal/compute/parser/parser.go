package parser

import (
	"github.com/DaniilZ77/InMemDB/internal/domain/models"
)

const (
	allowedSymbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

type tokenKind int

const (
	set tokenKind = iota
	get
	del
	argument
	eof
)

var reserved = map[string]tokenKind{
	"set": set,
	"get": get,
	"del": del,
}

type token struct {
	kind  tokenKind
	value string
}

func isIn(b byte, s string) bool {
	for i := range s {
		if b == s[i] {
			return true
		}
	}
	return false
}

func (p *Parser) tokenize(source string) []token {
	sourceLen := len(source)
	var pos, prevPos int
	var value string
	var tokens []token
	for {
		for ; pos < sourceLen && !isIn(source[pos], allowedSymbols); pos++ {
			prevPos++
		}

		pos++
		for ; pos < sourceLen && isIn(source[pos], allowedSymbols); pos++ {
		}

		if pos >= sourceLen {
			break
		}

		value = source[prevPos:pos]
		if kind, ok := reserved[value]; ok {
			tokens = append(tokens, token{
				kind:  kind,
				value: value,
			})
		} else {
			tokens = append(tokens, token{
				kind:  argument,
				value: value,
			})
		}

		prevPos = pos
	}

	tokens = append(tokens, token{
		kind:  eof,
		value: "EOF",
	})

	return tokens
}

func (p *Parser) Parse(source string) (*models.Command, error) {
	tokens := p.tokenize(source)

	var command models.Command
	if err := p.parseOperation(&command, tokens); err != nil {
		return nil, err
	}

	return &command, nil
}

func (p *Parser) parseOperation(command *models.Command, tokens []token) error {
	operation := tokens[0]
	switch operation.kind {
	case set:
		command.Type = models.SET
		return p.parseArgs(command, tokens[1:])
	case get:
		command.Type = models.GET
		return p.parseArg(command, tokens[1:])
	case del:
		command.Type = models.DEL
		return p.parseArg(command, tokens[1:])
	default:
		return models.ErrUnknownCommand
	}
}

func (p *Parser) parseArg(command *models.Command, tokens []token) error {
	arg := tokens[0]
	switch arg.kind {
	case argument:
		command.Args = append(command.Args, arg.value)
	default:
		return models.ErrUnknownCommand
	}

	return p.parseEof(tokens[1:])
}

func (p *Parser) parseArgs(command *models.Command, tokens []token) error {
	if len(tokens) < 2 {
		return models.ErrUnknownCommand
	}

	arg1, arg2 := tokens[0], tokens[1]
	if arg1.kind == argument && arg2.kind == argument {
		command.Args = append(command.Args, arg1.value, arg2.value)
		return p.parseExtraArgs(command, tokens[2:])
	} else {
		return models.ErrUnknownCommand
	}
}

func (p *Parser) parseExtraArgs(command *models.Command, tokens []token) error {
	if err := p.parseEof(tokens); err == nil {
		return nil
	}

	return p.parseArgs(command, tokens)
}

func (p *Parser) parseEof(tokens []token) error {
	if kind := tokens[0].kind; kind == eof {
		return nil
	}

	return models.ErrUnknownCommand
}

package parser

import "errors"

type CommandType int

const (
	GET CommandType = iota
	SET
	DEL

	setArgsCount     = 2
	defaultArgsCount = 1
)

var keywords = map[string]CommandType{
	"get": GET,
	"set": SET,
	"del": DEL,
}

type Command struct {
	Type CommandType
	Args []string
}

func (ct CommandType) argsCount() int {
	switch ct {
	case SET:
		return setArgsCount
	default:
		return defaultArgsCount
	}
}

var (
	ErrInvalidCommand = errors.New("invalid command")
)

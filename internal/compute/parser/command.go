package parser

import "errors"

type commandType int

const (
	GET commandType = iota
	SET
	DEL

	setArgsCount     = 2
	defaultArgsCount = 1
)

var keywords = map[string]commandType{
	"get": GET,
	"set": SET,
	"del": DEL,
}

type Command struct {
	Type commandType
	Args []string
}

var (
	ErrInvalidCommand = errors.New("invalid command")
)

func (ct commandType) argsCount() int {
	switch ct {
	case SET:
		return setArgsCount
	default:
		return defaultArgsCount
	}
}

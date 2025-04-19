package parser

import "errors"

type CommandType int

const (
	GET CommandType = iota
	SET
	DEL
	COMMIT
	ROLLBACK
	BEGIN

	setArgsCount     = 2
	getArgsCount     = 1
	delArgsCount     = 1
	defaultArgsCount = 0
)

var keywords = map[string]CommandType{
	"get":      GET,
	"set":      SET,
	"del":      DEL,
	"commit":   COMMIT,
	"rollback": ROLLBACK,
	"begin":    BEGIN,
}

type Command struct {
	Type CommandType
	Args []string
}

func NewDelCommand(key string) *Command {
	return &Command{
		Type: DEL,
		Args: []string{key},
	}
}

func NewGetCommand(key string) *Command {
	return &Command{
		Type: GET,
		Args: []string{key},
	}
}

func NewSetCommand(key, value string) *Command {
	return &Command{
		Type: SET,
		Args: []string{key, value},
	}
}

func NewBeginCommand() *Command {
	return &Command{Type: BEGIN}
}

func NewCommitCommand() *Command {
	return &Command{Type: COMMIT}
}

func NewRollbackCommand() *Command {
	return &Command{Type: ROLLBACK}
}

func (ct CommandType) argsCount() int {
	switch ct {
	case SET:
		return setArgsCount
	case GET:
		return getArgsCount
	case DEL:
		return delArgsCount
	case BEGIN, COMMIT, ROLLBACK:
		return defaultArgsCount
	default:
		panic("invalid command type")
	}
}

var (
	ErrInvalidCommand = errors.New("invalid command")
)

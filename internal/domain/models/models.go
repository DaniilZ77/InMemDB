package models

type CommandType int

const (
	SET CommandType = iota
	GET
	DEL
)

type Command struct {
	Type CommandType
	Args []string
}

type Pair struct {
	Key   string
	Value string
}

type SetCommand struct {
	Args      Pair
	ExtraArgs []Pair
}
type GetCommand string
type DelCommand string

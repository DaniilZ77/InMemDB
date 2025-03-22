package wal

import "github.com/DaniilZ77/InMemDB/internal/compute/parser"

type Command struct {
	CommandType int
	Args        []string
}

type batch struct {
	batchSize   int
	commands    []Command
	doneChannel chan bool
}

func newBatch(batchSize int) *batch {
	return &batch{batchSize: batchSize, doneChannel: make(chan bool)}
}

func (b *batch) appendCommand(command *parser.Command) {
	b.commands = append(b.commands, Command{
		CommandType: int(command.Type),
		Args:        command.Args,
	})
}

func (b *batch) resetBatch() {
	b.commands = nil
	b.doneChannel = make(chan bool)
}

func (b *batch) notifyFlushed(status bool) {
	defer close(b.doneChannel)
	for range b.commands {
		b.doneChannel <- status
	}
}

func (b *batch) isFull() bool {
	return len(b.commands) >= b.batchSize
}

func (b *batch) waitFlushed() bool {
	return <-b.doneChannel
}

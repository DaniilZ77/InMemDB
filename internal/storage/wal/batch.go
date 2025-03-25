package wal

import "github.com/DaniilZ77/InMemDB/internal/compute/parser"

type batch struct {
	lsn         int
	batchSize   int
	commands    []Command
	doneChannel chan bool
}

func NewBatch(batchSize int) *batch {
	return &batch{batchSize: batchSize, doneChannel: make(chan bool)}
}

func (b *batch) AppendCommand(command *parser.Command) {
	b.commands = append(b.commands, Command{
		LSN:         b.lsn,
		CommandType: int(command.Type),
		Args:        command.Args,
	})
	b.lsn++
}

func (b *batch) ResetBatch() {
	b.commands = nil
	b.doneChannel = make(chan bool)
}

func (b *batch) NotifyFlushed(status bool) {
	defer close(b.doneChannel)
	for range b.commands {
		b.doneChannel <- status
	}
}

func (b *batch) IsFull() bool {
	return len(b.commands) >= b.batchSize
}

func (b *batch) WaitFlushed() bool {
	return <-b.doneChannel
}

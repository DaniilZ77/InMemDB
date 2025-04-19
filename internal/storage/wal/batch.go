package wal

import "github.com/DaniilZ77/InMemDB/internal/compute/parser"

type Command struct {
	LSN         int
	TxID        int64
	CommandType int
	Args        []string
}

type Batch struct {
	lsn         int
	batchSize   int
	commands    []Command
	doneChannel chan bool
}

func NewBatch(batchSize int) *Batch {
	return &Batch{batchSize: batchSize, doneChannel: make(chan bool)}
}

func (b *Batch) AppendCommand(txID int64, command *parser.Command) {
	b.commands = append(b.commands, Command{
		LSN:         b.lsn,
		TxID:        txID,
		CommandType: int(command.Type),
		Args:        command.Args,
	})
	b.lsn++
}

func (b *Batch) ResetBatch() {
	b.commands = nil
	b.doneChannel = make(chan bool)
}

func (b *Batch) NotifyFlushed(status bool) {
	defer close(b.doneChannel)
	for range b.commands {
		b.doneChannel <- status
	}
}

func (b *Batch) IsFull() bool {
	return len(b.commands) >= b.batchSize
}

func (b *Batch) WaitFlushed() bool {
	return <-b.doneChannel
}

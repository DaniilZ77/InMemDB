package wal

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
)

const (
	statusSuccess = true
	statusError   = false
)

//go:generate mockery --name=LogsReader --case=snake --inpackage --inpackage-suffix --with-expecter
type LogsReader interface {
	ReadLogs() ([]Command, error)
}

//go:generate mockery --name=LogsWriter --case=snake --inpackage --inpackage-suffix --with-expecter
type LogsWriter interface {
	WriteLogs([]Command) error
}

type Wal struct {
	logsReader   LogsReader
	logsWriter   LogsWriter
	batchChannel chan Batch
	batchTimeout time.Duration
	log          *slog.Logger
	mu           sync.Mutex
	batch        *Batch
}

func NewWal(
	batchSize int,
	batchTimeout time.Duration,
	logsReader LogsReader,
	logsWriter LogsWriter,
	log *slog.Logger) (*Wal, error) {
	if logsReader == nil {
		return nil, errors.New("logs reader is nil")
	}
	if logsWriter == nil {
		return nil, errors.New("logs writer is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	return &Wal{
		logsReader:   logsReader,
		logsWriter:   logsWriter,
		batchChannel: make(chan Batch),
		batchTimeout: batchTimeout,
		log:          log,
		batch:        NewBatch(batchSize),
	}, nil
}

func (w *Wal) Save(txID int64, command *parser.Command) bool {
	w.mu.Lock()
	w.batch.AppendCommand(txID, command)
	batch := *w.batch
	if w.batch.IsFull() {
		w.batch.ResetBatch()
		w.batchChannel <- batch
	}
	w.mu.Unlock()

	return batch.WaitFlushed()
}

func (w *Wal) Start(ctx context.Context) {
	ticker := time.NewTicker(w.batchTimeout)

	defer func() {
		ticker.Stop()
		if v := recover(); v != nil {
			w.log.Error("panic recovered", slog.Any("error", v))
			go w.Start(ctx)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			w.log.Info("flushing and stopping wal")
			w.flushAll()
			return
		default:
		}
		select {
		case <-ctx.Done():
			w.log.Info("flushing and stopping wal")
			w.flushAll()
			return
		case <-ticker.C:
			w.mu.Lock()
			batch := *w.batch
			w.batch.ResetBatch()
			w.mu.Unlock()
			w.flushBatch(batch)
		case batch := <-w.batchChannel:
			ticker.Reset(w.batchTimeout)
			w.flushBatch(batch)
		}
	}
}

func (w *Wal) flushAll() {
	for {
		select {
		case batch := <-w.batchChannel:
			w.flushBatch(batch)
		default:
			w.mu.Lock()
			batch := *w.batch
			w.mu.Unlock()
			w.flushBatch(batch)
			return
		}
	}
}

func (w *Wal) flushBatch(batch Batch) {
	if len(batch.commands) == 0 {
		return
	}

	err := w.logsWriter.WriteLogs(batch.commands)
	if err != nil {
		w.log.Error("failed to flush batch", slog.Any("error", err))
		batch.NotifyFlushed(statusError)
		return
	}

	batch.NotifyFlushed(statusSuccess)
}

func (w *Wal) Recover() ([]Command, error) {
	commands, err := w.logsReader.ReadLogs()
	if err != nil {
		return nil, err
	}

	w.log.Info("recovered database", slog.Any("commands", len(commands)))
	if len(commands) > 0 {
		w.batch.lsn = commands[len(commands)-1].LSN + 1
	}

	return commands, nil
}

package wal

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

const (
	statusSuccess = true
	statusError   = false
)

type Wal struct {
	logsManager *logsManager

	batchChannel chan batch
	batchTimeout time.Duration

	log *slog.Logger

	mu    sync.Mutex
	batch *batch
}

func NewWal(cfg *config.Config, disk Disk, log *slog.Logger) (*Wal, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if disk == nil {
		return nil, errors.New("disk is nil")
	}
	if log == nil {
		return nil, errors.New("logger is nil")
	}

	return &Wal{
		logsManager:  NewLogsManager(disk, log),
		batchChannel: make(chan batch),
		batchTimeout: cfg.Wal.FlushingBatchTimeout,
		log:          log,
		batch:        NewBatch(cfg.Wal.FlushingBatchSize),
	}, nil
}

func (w *Wal) Save(command *parser.Command) bool {
	w.mu.Lock()
	w.batch.AppendCommand(command)
	batch := *w.batch
	if w.batch.IsFull() {
		w.batch.ResetBatch()
		w.mu.Unlock()
		w.batchChannel <- batch
	} else {
		w.mu.Unlock()
	}

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
			go w.flushBatch(batch)
		case batch := <-w.batchChannel:
			ticker.Reset(w.batchTimeout)
			go w.flushBatch(batch)
		}
	}
}

func (w *Wal) flushAll() {
	for {
		select {
		case batch := <-w.batchChannel:
			go w.flushBatch(batch)
		default:
			w.mu.Lock()
			batch := *w.batch
			w.mu.Unlock()
			go w.flushBatch(batch)
			return
		}
	}
}

func (w *Wal) flushBatch(batch batch) {
	if len(batch.commands) == 0 {
		return
	}

	err := w.logsManager.Write(batch.commands)
	if err != nil {
		w.log.Error("failed to flush batch", slog.Any("error", err))
		batch.NotifyFlushed(statusError)
		return
	}

	batch.NotifyFlushed(statusSuccess)
}

func (w *Wal) Recover() ([]parser.Command, error) {
	commands, err := w.logsManager.Read()
	if err != nil {
		return nil, err
	}

	w.log.Info("recovered database", slog.Any("commands", len(commands)))

	slices.SortFunc(commands, func(command1, command2 Command) int {
		return command1.LSN - command2.LSN
	})

	if len(commands) > 0 {
		w.batch.lsn = commands[len(commands)-1].LSN + 1
	}

	parserCommands := make([]parser.Command, len(commands))
	for i := range commands {
		parserCommands[i] = parser.Command{
			Type: parser.CommandType(commands[i].CommandType),
			Args: commands[i].Args,
		}
	}

	return parserCommands, nil
}

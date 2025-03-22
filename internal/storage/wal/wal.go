package wal

import (
	"errors"
	"log/slog"
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
		logsManager:  newLogsManager(disk, log),
		batchChannel: make(chan batch),
		batchTimeout: cfg.Wal.FlushingBatchTimeout,
		log:          log,
		batch:        newBatch(cfg.Wal.FlushingBatchSize),
	}, nil
}

func (w *Wal) Save(command *parser.Command) bool {
	w.mu.Lock()
	w.batch.appendCommand(command)
	batch := *w.batch
	if w.batch.isFull() {
		w.batch.resetBatch()
		w.mu.Unlock()
		w.batchChannel <- batch
	} else {
		w.mu.Unlock()
	}

	return batch.waitFlushed()
}

func (w *Wal) Start() {
	ticker := time.NewTicker(w.batchTimeout)

	defer func() {
		if v := recover(); v != nil {
			w.log.Error("panic recovered", slog.Any("error", v))
		}
		ticker.Stop()
		go w.Start()
	}()

	for {
		select {
		case <-ticker.C:
			w.mu.Lock()
			batch := *w.batch
			w.batch.resetBatch()
			w.mu.Unlock()
			go w.flushBatch(batch)
		case batch := <-w.batchChannel:
			ticker.Reset(w.batchTimeout)
			go w.flushBatch(batch)
		}
	}
}

func (w *Wal) flushBatch(batch batch) {
	if len(batch.commands) == 0 {
		return
	}

	err := w.logsManager.write(batch.commands)
	if err != nil {
		w.log.Error("failed to flush batch", slog.Any("error", err))
		batch.notifyFlushed(statusError)
		return
	}

	batch.notifyFlushed(statusSuccess)
}

func (w *Wal) Recover() ([]parser.Command, error) {
	commands, err := w.logsManager.read()
	if err != nil {
		return nil, err
	}

	w.log.Info("recovered database", slog.Any("commands", len(commands)))

	parserCommands := make([]parser.Command, len(commands))
	for i := range commands {
		parserCommands[i] = parser.Command{
			Type: parser.CommandType(commands[i].CommandType),
			Args: commands[i].Args,
		}
	}

	return parserCommands, nil
}

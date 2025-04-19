package wal

import (
	"bytes"
	"encoding/gob"
	"log/slog"

	"github.com/DaniilZ77/InMemDB/internal/common"
)

//go:generate mockery --name=Disk --case=snake --inpackage --inpackage-suffix --with-expecter
type Disk interface {
	WriteSegment([]byte) error
	ReadSegments() ([]byte, error)
}

type LogsManager struct {
	disk Disk
	log  *slog.Logger
}

func NewLogsManager(disk Disk, log *slog.Logger) *LogsManager {
	return &LogsManager{
		disk: disk,
		log:  log,
	}
}

func (w *LogsManager) WriteLogs(commands []Command) error {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(commands); err != nil {
		return err
	}

	if err := w.disk.WriteSegment(buffer.Bytes()); err != nil {
		w.log.Error("failed to write data on disk", slog.Any("error", err))
		return err
	}

	return nil
}

func (w *LogsManager) ReadLogs() ([]Command, error) {
	segments, err := w.disk.ReadSegments()
	if err != nil {
		return nil, err
	}

	return common.DecodeMany[[]Command](segments)
}

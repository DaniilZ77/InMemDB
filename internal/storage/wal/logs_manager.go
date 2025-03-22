package wal

import (
	"bytes"
	"io"
	"log/slog"
)

type Disk interface {
	Write([]byte) error
	Read() ([]byte, error)
}

type logsManager struct {
	disk Disk
	log  *slog.Logger
}

func newLogsManager(disk Disk, log *slog.Logger) *logsManager {
	return &logsManager{
		disk: disk,
		log:  log,
	}
}

func (w *logsManager) write(commands []Command) error {
	var buf []byte
	for _, command := range commands {
		encodedCommand, err := command.Encode()
		if err != nil {
			return err
		}
		buf = append(buf, encodedCommand...)
	}

	if err := w.disk.Write(buf); err != nil {
		w.log.Error("failed to write data on disk", slog.Any("error", err))
		return err
	}

	return nil
}

func (w *logsManager) read() ([]Command, error) {
	data, err := w.disk.Read()
	if err != nil {
		return nil, err
	}

	var commands []Command
	buf := bytes.NewBuffer(data)
	for {
		var command Command
		err = command.Decode(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

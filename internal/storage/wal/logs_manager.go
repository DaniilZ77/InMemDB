package wal

import (
	"bytes"
	"io"
	"log/slog"
)

//go:generate mockery --name=Disk --case=snake --inpackage --inpackage-suffix --with-expecter
type Disk interface {
	Write([]byte) error
	Read() ([]byte, error)
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

func (w *LogsManager) Write(commands []Command) error {
	var data []byte
	for _, command := range commands {
		encodedCommand, err := command.Encode()
		if err != nil {
			return err
		}
		data = append(data, encodedCommand...)
	}

	if err := w.disk.Write(data); err != nil {
		w.log.Error("failed to write data on disk", slog.Any("error", err))
		return err
	}

	return nil
}

func (w *LogsManager) Read() ([]Command, error) {
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

package wal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	for _, command := range commands {
		if err := encoder.Encode(command); err != nil {
			return err
		}
	}

	length := new(bytes.Buffer)
	if err := binary.Write(length, binary.LittleEndian, uint32(data.Len())); err != nil {
		return err
	}

	if err := w.disk.Write(append(length.Bytes(), data.Bytes()...)); err != nil {
		w.log.Error("failed to write data on disk", slog.Any("error", err))
		return err
	}

	return nil
}

func (w *logsManager) decodeCommands(data []byte) ([]Command, error) {
	var commands []Command
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	for {
		var command Command
		err := decoder.Decode(&command)
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

func (w *logsManager) decodeLength(buf *bytes.Buffer) (uint32, error) {
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return 0, err
	}

	return length, nil
}

func (w *logsManager) decodeBatch(buf *bytes.Buffer) ([]Command, error) {
	length, err := w.decodeLength(buf)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := buf.Read(data); err != nil {
		return nil, err
	}

	return w.decodeCommands(data)
}

func (w *logsManager) read() ([]Command, error) {
	data, err := w.disk.Read()
	if err != nil {
		return nil, err
	}

	var commands []Command
	buf := bytes.NewBuffer(data)
	for {
		res, err := w.decodeBatch(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		commands = append(commands, res...)
	}

	return commands, nil
}

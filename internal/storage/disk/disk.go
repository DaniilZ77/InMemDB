package disk

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var ErrSegmentNotFound = errors.New("segment not found")

type Disk struct {
	directory string
	segment   *Segment
	log       *slog.Logger
}

func NewDisk(dataDirectory string, maxSegmentSize int, log *slog.Logger) *Disk {
	return &Disk{
		directory: dataDirectory,
		segment:   NewSegment(maxSegmentSize, dataDirectory, log),
		log:       log,
	}
}

func (d *Disk) WriteSegment(data []byte) error {
	return d.segment.Write(data)
}

func (d *Disk) ReadSegments() ([]byte, error) {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return nil, err
	}

	var data []byte
	for i := range entries {
		segment, err := os.ReadFile(filepath.Join(d.directory, entries[i].Name()))
		if err != nil {
			return nil, err
		}

		data = append(data, segment...)
	}

	return data, nil
}

func (d *Disk) NextSegment(filename string) (string, error) {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", ErrSegmentNotFound
	}

	index, _ := slices.BinarySearchFunc(entries, filename, func(dir os.DirEntry, filename string) int {
		return strings.Compare(dir.Name(), filename)
	})

	if index >= len(entries)-1 {
		return entries[index].Name(), nil
	}

	return entries[index+1].Name(), nil
}

func (d *Disk) LastSegment() (string, error) {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", ErrSegmentNotFound
	}

	return entries[len(entries)-1].Name(), nil
}

func (d *Disk) WriteFile(filename string, data []byte) error {
	file, err := os.OpenFile(filepath.Join(d.directory, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			d.log.Warn("failed to close file", slog.String("filename", filename), slog.Any("error", err))
		}
	}()

	if _, err = file.Write(data); err != nil {
		return err
	}

	return file.Sync()
}

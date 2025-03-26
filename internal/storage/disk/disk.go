package disk

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Disk struct {
	directory string
	mu        sync.RWMutex
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

func (d *Disk) Write(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.segment.Write(data)
}

func (d *Disk) read(fileName string) (data []byte, err error) {
	file, err := os.Open(filepath.Join(d.directory, fileName))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			d.log.Error("failed to close file", slog.String("filename", fileName), slog.Any("error", err))
		}
	}()

	return io.ReadAll(file)
}

func (d *Disk) Read() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return nil, err
	}

	var data []byte
	for i := range entries {
		if entries[i].IsDir() {
			continue
		}

		segment, err := d.read(entries[i].Name())
		if err != nil {
			return nil, err
		}

		data = append(data, segment...)
	}

	return data, nil
}

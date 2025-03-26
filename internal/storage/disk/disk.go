package disk

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/DaniilZ77/InMemDB/internal/config"
)

type Disk struct {
	directory string
	mu        sync.RWMutex
	segment   *Segment
	log       *slog.Logger
}

func NewDisk(cfg *config.Config, log *slog.Logger) (*Disk, error) {
	return &Disk{
		directory: cfg.Wal.DataDirectory,
		segment:   NewSegment(cfg.Wal.MaxSegmentSize, cfg.Wal.DataDirectory, log),
		log:       log,
	}, nil
}

func (d *Disk) Write(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.segment.Write(data); err != nil {
		return err
	}

	return nil
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

	data, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
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

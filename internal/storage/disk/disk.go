package disk

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/DaniilZ77/InMemDB/internal/config"
)

type Disk struct {
	directory string
	mu        sync.RWMutex
	segment   *segment
}

func NewDisk(cfg *config.Config, log *slog.Logger) (*Disk, error) {
	return &Disk{
		directory: cfg.Wal.DataDirectory,
		segment:   newSegment(cfg.Wal.MaxSegmentSize, 0, cfg.Wal.DataDirectory),
	}, nil
}

func (d *Disk) Write(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.segment.write(data); err != nil {
		return err
	}

	return nil
}

func (d *Disk) Read() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		if err := d.segment.newFile(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	buf := new(bytes.Buffer)
	for i := range entries {
		if entries[i].IsDir() {
			continue
		}

		path := filepath.Join(d.directory, entries[i].Name())
		data, err := d.segment.read(path)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(data); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

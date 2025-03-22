package disk

import (
	"log/slog"
	"os"
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
		segment:   newSegment(cfg.Wal.MaxSegmentSize, cfg.Wal.DataDirectory),
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

	var lastSegmentIndex int
	segments := make([][]byte, len(entries))
	for i := range entries {
		if entries[i].IsDir() {
			continue
		}

		data, segmentIndex, err := d.segment.read(entries[i].Name())
		if err != nil {
			return nil, err
		}
		segments[segmentIndex] = data
		lastSegmentIndex = max(lastSegmentIndex, segmentIndex)
	}

	d.segment.curSegmentIndex = lastSegmentIndex
	if err := d.segment.newFile(); err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		d.segment.curSegmentSize = len(segments[lastSegmentIndex])
	}

	var data []byte
	for i := range segments {
		data = append(data, segments[i]...)
	}

	return data, nil
}

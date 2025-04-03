package disk

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type Disk struct {
	directory string
	segment   *Segment
	log       *slog.Logger
}

func NewDisk(dataDirectory string, maxSegmentSize int, log *slog.Logger) *Disk {
	if err := os.MkdirAll(dataDirectory, 0777); err != nil && !errors.Is(err, fs.ErrExist) {
		log.Error("failed to create directory", slog.String("directory", dataDirectory), slog.Any("error", err))
		return nil
	}

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

	return upperBound(entries, filename), nil
}

func upperBound(entries []os.DirEntry, filename string) string {
	l, r := 0, len(entries)-1
	for l <= r {
		m := (l + r) / 2
		if entries[m].Name() <= filename {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	if l >= len(entries)-1 {
		return ""
	}
	return entries[l].Name()
}

func (d *Disk) LastSegment() (string, error) {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", nil
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

package disk

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type segment struct {
	maxSegmentSize int
	curSegmentSize int
	directory      string
	file           *os.File
	log            *slog.Logger
}

func NewSegment(maxSegmentSize int, directory string, log *slog.Logger) *segment {
	return &segment{
		maxSegmentSize: maxSegmentSize,
		directory:      directory,
		log:            log,
	}
}

func (s *segment) rotateSegment() (err error) {
	if s.file != nil {
		if err := s.file.Close(); err != nil {
			s.log.Error("failed to close file", slog.String("filename", s.file.Name()), slog.Any("error", err))
		}
	}

	fileName := fmt.Sprintf("wal_%d.log", time.Now().UnixMilli())
	s.file, err = os.OpenFile(filepath.Join(s.directory, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	s.curSegmentSize = 0
	return nil
}

func (s *segment) Write(data []byte) error {
	if s.file == nil || s.curSegmentSize >= s.maxSegmentSize {
		if err := s.rotateSegment(); err != nil {
			return err
		}
	}

	written, err := s.file.Write(data)
	if err != nil {
		return err
	}
	if err := s.file.Sync(); err != nil {
		return err
	}

	s.curSegmentSize += written
	return nil
}

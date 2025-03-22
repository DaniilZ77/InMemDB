package disk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type segment struct {
	maxSegmentSize  int
	curSegmentSize  int
	curSegmentIndex int
	directory       string
	file            *os.File
}

func newSegment(maxSegmentSize int, directory string) *segment {
	return &segment{
		maxSegmentSize: maxSegmentSize,
		directory:      directory,
	}
}

func (s *segment) newFile() (err error) {
	if s.file != nil {
		s.file.Close()
	}

	fileName := fmt.Sprintf("wal_%d.log", s.curSegmentIndex)
	s.file, err = os.OpenFile(filepath.Join(s.directory, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	s.curSegmentIndex++
	s.curSegmentSize = 0
	return nil
}

func (s *segment) write(data []byte) error {
	if s.curSegmentSize >= s.maxSegmentSize {
		if err := s.newFile(); err != nil {
			return err
		}
	}

	if _, err := s.file.Write(data); err != nil {
		return err
	}
	if err := s.file.Sync(); err != nil {
		return err
	}
	s.curSegmentSize += len(data)
	return nil
}

func (s *segment) read(fileName string) (data []byte, segmentIndex int, err error) {
	if _, err := fmt.Sscanf(fileName, "wal_%d.log", &segmentIndex); err != nil {
		return nil, 0, fmt.Errorf("failed to parse file name (%s) of wal log: %w", fileName, err)
	}

	file, err := os.Open(filepath.Join(s.directory, fileName))
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	data, err = io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	return data, segmentIndex, nil
}

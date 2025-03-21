package disk

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type segment struct {
	maxSegmentSize int
	curSegmentSize int
	directory      string
	file           *os.File
	modifyTime     time.Time
}

func newSegment(maxSegmentSize int, curSegmentSize int, directory string) *segment {
	return &segment{
		maxSegmentSize: maxSegmentSize,
		curSegmentSize: curSegmentSize,
		directory:      directory,
	}
}

func (s *segment) newFile() (err error) {
	if s.file != nil {
		s.file.Close()
	}

	fileName := fmt.Sprintf("wal_%s.log", uuid.NewString())
	s.file, err = os.OpenFile(filepath.Join(s.directory, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

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

func (s *segment) read(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	modifyTime := stat.ModTime()
	if modifyTime.After(s.modifyTime) {
		if s.file != nil {
			s.file.Close()
		}
		s.modifyTime = modifyTime
		s.file = file
	} else {
		defer file.Close()
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	s.curSegmentSize = int(stat.Size())
	return buf.Bytes(), nil
}

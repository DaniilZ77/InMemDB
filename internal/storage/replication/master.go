package replication

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

type NextSegmentProvider interface {
	NextSegment(filename string) (string, error)
}

type Master struct {
	disk         NextSegmentProvider
	walDirectory string
	log          *slog.Logger
}

func NewMaster(disk NextSegmentProvider, walDirectory string, log *slog.Logger) *Master {
	return &Master{
		disk:         disk,
		walDirectory: walDirectory,
		log:          log,
	}
}

func (m *Master) IsSlave() bool {
	return false
}

func (m *Master) GetReplicationStream() <-chan []wal.Command {
	return nil
}

func (m *Master) HandleRequest(request []byte) ([]byte, error) {
	decodedRequest, err := common.DecodeOne[Request](request)
	if err != nil {
		return nil, err
	}

	filename, err := m.disk.NextSegment(decodedRequest.LastSegment)
	if err != nil {
		return nil, err
	}

	if filename == "" {
		return common.Encode(NewResponse("", nil))
	}

	segment, err := os.ReadFile(filepath.Join(m.walDirectory, filename))
	if err != nil {
		return nil, err
	}

	return common.Encode(NewResponse(filename, segment))
}

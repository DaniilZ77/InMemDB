package replication

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

//go:generate mockery --name=NextSegmentProvider --case=snake --inpackage --inpackage-suffix --with-expecter
type NextSegmentProvider interface {
	NextSegment(filename string) (string, error)
}

type Master struct {
	disk         NextSegmentProvider
	walDirectory string
	log          *slog.Logger
}

func NewMaster(disk NextSegmentProvider, walDirectory string, log *slog.Logger) (*Master, error) {
	if disk == nil {
		return nil, errors.New("disk is nil")
	}
	if log == nil {
		return nil, errors.New("log is nil")
	}

	return &Master{
		disk:         disk,
		walDirectory: walDirectory,
		log:          log,
	}, nil
}

func (m *Master) IsSlave() bool {
	return false
}

func (m *Master) GetReplicationStream() <-chan []wal.Command {
	return nil
}

func (m *Master) HandleRequest(request []byte) (response []byte, err error) {
	defer func() {
		if err != nil {
			m.log.Warn("failed to handle request", slog.Any("error", err))
			response, err = common.Encode(NewErrorResponse())
		}
	}()

	var decodedRequest Request
	decodedRequest, err = common.DecodeOne[Request](request)
	if err != nil {
		return
	}

	m.log.Debug("received request from slave", slog.String("last_segment", decodedRequest.LastSegment))

	var filename string
	filename, err = m.disk.NextSegment(decodedRequest.LastSegment)
	if err != nil {
		return
	}

	var segment []byte
	segment, err = os.ReadFile(filepath.Join(m.walDirectory, filename))
	if err != nil {
		return
	}

	return common.Encode(NewSuccessResponse(filename, segment))
}

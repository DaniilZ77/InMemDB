package replication

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
)

//go:generate mockery --name=SegmentManager --case=snake --inpackage --inpackage-suffix --with-expecter
type SegmentManager interface {
	LastSegment() (string, error)
	WriteFile(filename string, data []byte) error
}

//go:generate mockery --name=Client --case=snake --inpackage --inpackage-suffix --with-expecter
type Client interface {
	Send(request []byte) ([]byte, error)
	Close() error
}

type Slave struct {
	syncInterval      time.Duration
	lastSegment       string
	replicationStream chan []wal.Command
	client            Client
	segmentManager    SegmentManager
	log               *slog.Logger
}

func NewSlave(
	syncInterval time.Duration,
	client Client,
	segmentManager SegmentManager,
	log *slog.Logger) (*Slave, error) {
	if segmentManager == nil {
		return nil, errors.New("segment manager is nil")
	}
	if log == nil {
		return nil, errors.New("log is nil")
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}

	lastSegment, err := segmentManager.LastSegment()
	if err != nil {
		return nil, err
	}

	return &Slave{
		syncInterval:      syncInterval,
		lastSegment:       lastSegment,
		replicationStream: make(chan []wal.Command),
		client:            client,
		segmentManager:    segmentManager,
		log:               log,
	}, nil
}

func (s *Slave) GetReplicationStream() <-chan []wal.Command {
	return s.replicationStream
}

func (s *Slave) Start(ctx context.Context) (err error) {
	ticker := time.NewTicker(s.syncInterval)

	defer func() {
		ticker.Stop()
		close(s.replicationStream)
		if err := s.client.Close(); err != nil {
			s.log.Warn("failed to close client", slog.Any("error", err))
		}
		if v := recover(); v != nil {
			s.log.Error("panic recovered", slog.Any("error", v))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("stopping slave")
			return nil
		default:
		}

		select {
		case <-ctx.Done():
			s.log.Info("stopping slave")
			return nil
		case <-ticker.C:
			if err := s.handle(ctx); err != nil {
				s.log.Error("failed to handle replication", slog.Any("error", err))
			}
		}
	}
}

func (s *Slave) IsSlave() bool {
	return true
}

func (s *Slave) receiveSegment() (*Response, error) {
	request := NewRequest(s.lastSegment)
	encodedRequest, err := common.Encode(request)
	if err != nil {
		return nil, err
	}

	response, err := s.client.Send(encodedRequest)
	if err != nil {
		s.log.Warn("failed to receive segment from master", slog.Any("error", err))
		return nil, err
	}

	decodedResponse, err := common.DecodeOne[Response](response)
	if err != nil {
		return nil, err
	}

	return &decodedResponse, nil
}

func (s *Slave) handle(ctx context.Context) error {
	response, err := s.receiveSegment()
	if err != nil {
		return err
	}

	s.log.Debug("received response from master",
		slog.Bool("ok", response.Ok),
		slog.String("filename", response.Filename),
	)

	if !response.Ok {
		return nil
	}

	if response.Filename == "" {
		s.log.Debug("no new segments")
		return nil
	}

	err = s.segmentManager.WriteFile(response.Filename, response.Segment)
	if err != nil {
		return err
	}

	decodedData, err := common.DecodeMany[[]wal.Command](response.Segment)
	if err != nil {
		return err
	}

	s.lastSegment = response.Filename

	select {
	case s.replicationStream <- decodedData:
	case <-ctx.Done():
	}
	return nil
}

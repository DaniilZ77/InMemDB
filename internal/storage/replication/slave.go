package replication

import (
	"context"
	"log/slog"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
	"github.com/DaniilZ77/InMemDB/internal/tcp/client"
)

type Disk interface {
	LastSegment() (string, error)
	WriteFile(filename string, data []byte) error
}

type Slave struct {
	masterAddress     string
	syncInterval      time.Duration
	bufferSize        int
	walDirectory      string
	lastSegment       string
	replicationStream chan []wal.Command
	client            *client.Client
	disk              Disk
	log               *slog.Logger
}

func NewSlave(
	masterAddress string,
	syncInterval time.Duration,
	bufferSize int,
	walDirectory string,
	disk Disk,
	log *slog.Logger) (*Slave, error) {
	lastSegment, err := disk.LastSegment()
	if err != nil {
		return nil, err
	}

	return &Slave{
		masterAddress:     masterAddress,
		syncInterval:      syncInterval,
		bufferSize:        bufferSize,
		walDirectory:      walDirectory,
		lastSegment:       lastSegment,
		replicationStream: make(chan []wal.Command),
		disk:              disk,
		log:               log,
	}, nil
}

func (s *Slave) GetReplicationStream() <-chan []wal.Command {
	return s.replicationStream
}

func (s *Slave) Start(ctx context.Context) (err error) {
	ticker := time.NewTicker(s.syncInterval)

	s.client, err = client.NewClient(s.masterAddress, s.bufferSize)
	if err != nil {
		return err
	}

	defer func() {
		ticker.Stop()
		close(s.replicationStream)
		if err := s.client.Close(); err != nil {
			s.log.Error("failed to close client", slog.Any("error", err))
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
				return err
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
		return nil, err
	}

	decodedResponse, err := common.DecodeOne[Response](response)
	if err != nil {
		return nil, err
	}

	s.lastSegment = decodedResponse.Filename
	return &decodedResponse, nil
}

func (s *Slave) handle(ctx context.Context) error {
	response, err := s.receiveSegment()
	if err != nil {
		return err
	}

	if !response.Ok {
		s.log.Warn("error response from master")
		return nil
	}

	err = s.disk.WriteFile(response.Filename, response.Segment)
	if err != nil {
		return err
	}

	decodedData, err := common.DecodeMany[[]wal.Command](response.Segment)
	if err != nil {
		return err
	}

	select {
	case s.replicationStream <- decodedData:
	case <-ctx.Done():
		s.log.Info("stopping slave")
	}
	return nil
}

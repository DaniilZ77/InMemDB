package app

import (
	"log/slog"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/storage/disk"
	"github.com/DaniilZ77/InMemDB/internal/storage/replication"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
	"github.com/DaniilZ77/InMemDB/internal/tcp/client"
)

const (
	defaultMaxSegmentSize       = 10 << 20
	defaultFlushingBatchTimeout = 10 * time.Millisecond
	defaultFlushingBatchSize    = 100
	defaultDataDirectory        = "./data/wal"
	slave                       = "slave"
	master                      = "master"
	defaultReplicaType          = master
	defaultMasterAddress        = ":3232"
	defaultSyncInterval         = time.Second
	defaultIdleTimeout          = time.Second
)

var replicaTypes = map[string]bool{
	slave:  true,
	master: true,
}

func NewWalReplica(config *config.Config, log *slog.Logger) (*wal.Wal, any, error) {
	if config.Wal == nil {
		return nil, nil, nil
	}

	replicaType := defaultReplicaType
	masterAddress := defaultMasterAddress
	syncInterval := defaultSyncInterval
	if config.Replication != nil {
		if replicaTypes[config.Replication.ReplicaType] {
			replicaType = config.Replication.ReplicaType
		}
		if config.Replication.MasterAddress != "" {
			masterAddress = config.Replication.MasterAddress
		}
		if config.Replication.SyncInterval > 0 {
			syncInterval = config.Replication.SyncInterval
		}
	}

	flushingBatchTimeout := defaultFlushingBatchTimeout
	flushingBatchSize := defaultFlushingBatchSize
	dataDirectory := defaultDataDirectory
	if config.Wal.FlushingBatchTimeout > 0 {
		flushingBatchTimeout = config.Wal.FlushingBatchTimeout
	}
	if config.Wal.FlushingBatchSize > 0 {
		flushingBatchSize = config.Wal.FlushingBatchSize
	}
	if config.Wal.DataDirectory != "" {
		dataDirectory = config.Wal.DataDirectory
	}

	maxSegmentSize, err := parseBytes(config.Wal.MaxSegmentSize)
	if err != nil {
		maxSegmentSize = defaultMaxSegmentSize
	}

	disk := disk.NewDisk(dataDirectory, maxSegmentSize, log)
	logsManager := wal.NewLogsManager(disk, log)

	wal, err := wal.NewWal(flushingBatchSize, flushingBatchTimeout, logsManager, logsManager, log)
	if err != nil {
		return nil, nil, err
	}

	if config.Replication == nil {
		return wal, nil, nil
	}

	var replica any
	switch replicaType {
	case master:
		replica, err = replication.NewMaster(disk, dataDirectory, log)
		return wal, replica, err
	case slave:
		client, err := NewClient(
			masterAddress,
			log,
			client.WithBufferSize(2*maxSegmentSize),
			client.WithIdleTimeout(defaultIdleTimeout),
		)
		if err != nil {
			return nil, nil, err
		}
		replica, err := replication.NewSlave(syncInterval, client, disk, log)
		return wal, replica, err
	}

	panic("unreachable")
}

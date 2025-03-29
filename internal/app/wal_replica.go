package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/storage/disk"
	"github.com/DaniilZ77/InMemDB/internal/storage/replication"
	"github.com/DaniilZ77/InMemDB/internal/storage/wal"
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
)

var replicaTypes = map[string]bool{
	slave:  true,
	master: true,
}

func NewWalReplica(ctx context.Context, config *config.Config, log *slog.Logger) (*wal.Wal, any, error) {
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

	flushingBatchTimeout := config.Wal.FlushingBatchTimeout
	flushingBatchSize := config.Wal.FlushingBatchSize
	dataDirectory := config.Wal.DataDirectory
	if flushingBatchTimeout <= 0 {
		flushingBatchSize = defaultFlushingBatchSize
	}
	if flushingBatchSize <= 0 {
		flushingBatchSize = defaultFlushingBatchSize
	}
	if dataDirectory == "" {
		dataDirectory = defaultDataDirectory
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
	if replicaType != slave {
		go wal.Start(ctx)
	}

	if config.Replication == nil {
		return wal, nil, nil
	}

	switch replicaType {
	case master:
		return wal, replication.NewMaster(disk, dataDirectory, log), nil
	case slave:
		replica, err := replication.NewSlave(masterAddress, syncInterval, defaultMaxSegmentSize, dataDirectory, disk, log)
		return wal, replica, err
	}

	panic("unreachable")
}

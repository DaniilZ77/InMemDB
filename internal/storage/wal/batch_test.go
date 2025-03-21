package wal

import (
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/compute/parser"
	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/storage/disk"
)

func Test(t *testing.T) {
	cfg := &config.Config{
		Wal: &config.Wal{
			FlushingBatchTimeout: 10 * time.Millisecond,
			FlushingBatchSize:    100,
			MaxSegmentSize:       10_000_000,
			DataDirectory:        "./../../../data/wal",
		},
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	disk, err := disk.NewDisk(cfg, log)
	if err != nil {
		panic(err)
	}

	wal, err := NewWal(cfg, disk, log)
	if err != nil {
		panic(err)
	}

	go wal.Start()

	wal.Recover()

	wg := sync.WaitGroup{}
	wg.Add(1000)

	for range 1000 {
		go func() {
			defer wg.Done()
			for range 10 {
				wal.Save(&parser.Command{Type: 1, Args: []string{"key", "value"}})
			}
		}()
	}

	wg.Wait()
}

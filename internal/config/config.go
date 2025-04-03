package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine      *Engine      `yaml:"engine"`
	Network     *Network     `yaml:"network"`
	LogLevel    string       `yaml:"log_level"`
	Wal         *Wal         `yaml:"wal"`
	Replication *Replication `yaml:"replication"`
}

type Network struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type Engine struct {
	Type         string `yaml:"type"`
	ShardsNumber int    `yaml:"shards_number"`
}

type Wal struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

type Replication struct {
	ReplicaType   string        `yaml:"replica_type"`
	MasterAddress string        `yaml:"master_address"`
	SyncInterval  time.Duration `yaml:"sync_interval"`
}

func MustConfig() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func NewConfig() (*Config, error) {
	configPath, ok := getConfigPath()
	if !ok {
		return nil, errors.New("config path is not set")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

func getConfigPath() (configPath string, ok bool) {
	flag.StringVar(&configPath, "config_path", "", "path to config")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath, configPath != ""
}

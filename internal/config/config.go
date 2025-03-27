package config

import (
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Engine   Engine  `yaml:"engine"`
	Network  Network `yaml:"network"`
	LogLevel string  `yaml:"log_level" env-default:"info"`
	Wal      *Wal    `yaml:"wal"`
}

type Network struct {
	Address             string        `yaml:"address" env-default:"127.0.0.1:3223"`
	MaxConnections      int           `yaml:"max_connections" env-default:"100"`
	MaxMessageSizeBytes int           `yaml:"-"`
	MaxMessageSize      string        `yaml:"max_message_size" env-default:"4KB"`
	IdleTimeout         time.Duration `yaml:"idle_timeout" env-default:"5m"`
}

type Engine struct {
	Type         string `yaml:"type" env-default:"in_memory"`
	ShardsNumber int    `yaml:"shards_number" env-default:"16"`
}

type Wal struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size" env-default:"100"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout" env-default:"10ms"`
	MaxSegmentSizeBytes  int           `yaml:"-"`
	MaxSegmentSize       string        `yaml:"max_segment_size" env-default:"10MB"`
	DataDirectory        string        `yaml:"data_directory" env-default:"/data/wal"`
}

func NewConfig() *Config {
	configPath, ok := getConfigPath()
	if !ok {
		panic("config path is not set")
	}

	var cfg Config
	var err error

	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	cfg.Network.MaxMessageSizeBytes, err = parseSize(cfg.Network.MaxMessageSize)
	if err != nil {
		panic("failed to parse max message size: " + err.Error())
	}

	cfg.Wal.MaxSegmentSizeBytes, err = parseSize(cfg.Wal.MaxSegmentSize)
	if err != nil {
		panic("failed to parse max segment size: " + err.Error())
	}

	return &cfg
}

func getConfigPath() (configPath string, ok bool) {
	flag.StringVar(&configPath, "config_path", "", "path to config")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath, configPath != ""
}

func parseSize(size string) (int, error) {
	sizeBytes := 0
	var unitMeasure string
	var i int
	for ; i < len(size) && '0' <= size[i] && size[i] <= '9'; i++ {
		sizeBytes = sizeBytes*10 + int(size[i]-'0')
	}
	unitMeasure = strings.TrimSpace(size[i:])

	switch strings.ToUpper(unitMeasure) {
	case "B":
		return sizeBytes, nil
	case "KB":
		return sizeBytes << 10, nil
	case "MB":
		return sizeBytes << 20, nil
	case "GB":
		return sizeBytes << 30, nil
	default:
		return 0, errors.New("invalid unit of measure")
	}
}

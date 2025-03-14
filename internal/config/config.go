package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Engine   Engine  `yaml:"engine"`
	Network  Network `yaml:"network"`
	LogLevel string  `yaml:"log_level" env-default:"info"`
}

type Network struct {
	Address        string        `yaml:"address" env-default:"127.0.0.1:3223"`
	MaxConnections int           `yaml:"max_connections" env-default:"100"`
	MaxMessageSize int           `yaml:"max_message_size" env-default:"4000"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"5m"`
}

type Engine struct {
	Type         string `yaml:"type" env-default:"in_memory"`
	ShardsAmount int    `yaml:"shards_amount" env-default:"10"`
}

func NewConfig() *Config {
	configPath, ok := getConfigPath()
	if !ok {
		panic("config path is not set")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
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

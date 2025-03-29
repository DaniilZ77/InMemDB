package app

import (
	"errors"

	"github.com/DaniilZ77/InMemDB/internal/config"
	"github.com/DaniilZ77/InMemDB/internal/storage/engine"
)

const (
	inMemory = "in_memory"

	defaultShardsNumber = 16
	defaultEngineType   = inMemory
)

var engineTypes = map[string]bool{
	inMemory: true,
}

func NewEngine(config *config.Config) (*engine.Engine, error) {
	shardsNumber := defaultShardsNumber

	if config.Engine != nil {
		if !engineTypes[config.Engine.Type] {
			return nil, errors.New("invalid engine type")
		}
		if config.Engine.ShardsNumber > 0 {
			shardsNumber = config.Engine.ShardsNumber
		}
	}

	return engine.NewEngine(shardsNumber)
}

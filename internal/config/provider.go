package config

import (
	"errors"

	"github.com/Fox1N69/iq-testtask/pkg/logger"
	"github.com/google/wire"
)

func ProvideConfig() (*Config, error) {
	cfg := LoadConfig("config/.env")
	if cfg == nil {
		return nil, errors.New("failed to load config")
	}

	logger.Init(cfg.Env.Mode)
	return cfg, nil
}

var ProviderSet = wire.NewSet(ProvideConfig)

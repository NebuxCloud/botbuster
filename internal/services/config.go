package services

import (
	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/samber/do/v2"
)

func NewConfig(i do.Injector) (*config.Config, error) {
	cfg := &config.Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

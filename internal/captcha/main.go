package captcha

import (
	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/NebuxCloud/botbuster/internal/data"
)

type Manager struct {
	cfg  *config.Config
	data data.Storage
}

func New(cfg *config.Config, d data.Storage) *Manager {
	return &Manager{
		cfg:  cfg,
		data: d,
	}
}

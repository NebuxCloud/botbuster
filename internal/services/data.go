package services

import (
	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/NebuxCloud/botbuster/internal/data"
	"github.com/samber/do/v2"
)

func NewData(i do.Injector) (data.Storage, error) {
	cfg := do.MustInvoke[*config.Config](i)

	return data.NewValkeyStorage(cfg)
}

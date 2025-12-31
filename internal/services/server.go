package services

import (
	"log/slog"

	"github.com/NebuxCloud/botbuster/internal/api"
	"github.com/NebuxCloud/botbuster/internal/captcha"
	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/samber/do/v2"
)

func NewServer(i do.Injector) (*api.API, error) {
	return api.New(
		do.MustInvoke[*config.Config](i),
		do.MustInvoke[*slog.Logger](i),
		do.MustInvoke[*captcha.Manager](i),
	), nil
}

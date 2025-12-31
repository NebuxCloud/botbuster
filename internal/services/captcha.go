package services

import (
	"github.com/NebuxCloud/botbuster/internal/captcha"
	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/NebuxCloud/botbuster/internal/data"
	"github.com/samber/do/v2"
)

func NewCaptcha(i do.Injector) (*captcha.Manager, error) {
	return captcha.New(
		do.MustInvoke[*config.Config](i),
		do.MustInvoke[data.Storage](i),
	), nil
}

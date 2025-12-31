package services

import (
	"log/slog"
	"os"

	"github.com/NebuxCloud/botbuster/internal/config"
	"github.com/samber/do/v2"
)

func NewLogger(i do.Injector) (*slog.Logger, error) {
	cfg := do.MustInvoke[*config.Config](i)

	level := slog.LevelInfo

	if cfg.Debug {
		level = slog.LevelDebug
	}

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	return log, nil
}

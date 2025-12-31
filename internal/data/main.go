package data

import (
	"context"
	"time"

	"github.com/samber/do/v2"
)

type Storage interface {
	do.ShutdownerWithContextAndError

	MarkUsed(ctx context.Context, id string, ttl time.Duration) error
	IsUsed(ctx context.Context, id string) (bool, error)
}

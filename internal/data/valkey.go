package data

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/NebuxCloud/botbuster/internal/config"
)

type ValkeyStorage struct {
	client valkey.Client
	prefix string
}

func NewValkeyStorage(cfg *config.Config) (Storage, error) {
	opt, err := valkey.ParseURL(cfg.ValkeyURL)

	if err != nil {
		return nil, err
	}

	client, err := valkey.NewClient(opt)

	if err != nil {
		return nil, err
	}

	return &ValkeyStorage{
		client: client,
		prefix: cfg.ValkeyPrefix,
	}, nil
}

func (s *ValkeyStorage) Shutdown(ctx context.Context) error {
	s.client.Close()
	return nil
}

func (s *ValkeyStorage) MarkUsed(ctx context.Context, id string, ttl time.Duration) error {
	res := s.client.Do(ctx, s.client.B().
		Set().
		Key(s.key(id)).
		Value("1").
		Ex(ttl).
		Build(),
	)

	return res.Error()
}

func (s *ValkeyStorage) IsUsed(ctx context.Context, id string) (bool, error) {
	res := s.client.Do(ctx, s.client.B().
		Exists().
		Key(s.key(id)).
		Build(),
	)

	if err := res.Error(); err != nil {
		return false, err
	}

	exists, err := res.ToInt64()

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (s *ValkeyStorage) key(id string) string {
	return s.prefix + "used:" + id
}

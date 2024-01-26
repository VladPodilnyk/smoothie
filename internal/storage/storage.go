package storage

import (
	"context"
)

type Storage interface {
	Inc(ctx context.Context, key string)
	Get(ctx context.Context, key string) (uint, error)
}

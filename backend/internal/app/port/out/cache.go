package out

import (
	"context"
	"time"
)

type CachePort interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	Evict(ctx context.Context, pattern string) error
}

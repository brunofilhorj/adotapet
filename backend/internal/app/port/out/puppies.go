package out

import (
	"context"

	"adotapet/internal/domain/common"
	"adotapet/internal/domain/puppy"
)

type PuppyRepository interface {
	Save(ctx context.Context, puppy puppy.Puppy) (puppy.Puppy, error)
	FindByID(ctx context.Context, id string) (*puppy.Puppy, error)
	FindByOwnerID(ctx context.Context, ownerID string, page common.PageRequest) (common.Page[puppy.Puppy], error)
}

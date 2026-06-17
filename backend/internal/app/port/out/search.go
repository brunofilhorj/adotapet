package out

import (
	"context"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/common"
)

type SearchRepository interface {
	SearchNearby(ctx context.Context, query inport.PuppySearchQuery) (common.Page[inport.PuppySummary], error)
}

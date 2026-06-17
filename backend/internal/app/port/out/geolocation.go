package out

import (
	"context"

	"adotapet/internal/domain/common"
)

type GeoLocationPort interface {
	Geocode(ctx context.Context, address Address) (common.GeoPoint, error)
}

type Address struct {
	City   string
	State  string
	Street string
	Number string
}

package in

import (
	"context"

	"adotapet/internal/domain/common"
)

type PuppySearchQuery struct {
	Latitude     float64
	Longitude    float64
	RadiusKM     float64
	Species      *string
	Breed        *string
	AgeMinMonths *int
	AgeMaxMonths *int
	Size         *string
	Sex          *string
	Page         common.PageRequest
}

type PuppySummary struct {
	ID              string
	Name            string
	Breed           *string
	Species         string
	AgeMonths       int
	Size            string
	Sex             string
	Status          string
	DistanceKM      float64
	City            string
	State           string
	PrimaryPhotoURL *string
}

type SearchPuppiesInputPort interface {
	Search(ctx context.Context, query PuppySearchQuery) (common.Page[PuppySummary], error)
}

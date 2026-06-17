package in

import (
	"context"

	"adotapet/internal/domain/common"
)

type CreatePuppyCommand struct {
	OwnerID     string
	OwnerRole   string
	Name        string
	Breed       *string
	Species     string
	AgeMonths   int
	Size        string
	Sex         string
	Description string
	Location    common.GeoPoint
	City        string
	State       string
}

type PuppyDetails struct {
	ID          string
	OwnerID     string
	Name        string
	Breed       *string
	Species     string
	AgeMonths   int
	Size        string
	Sex         string
	Description string
	Location    common.GeoPoint
	City        string
	State       string
	Status      string
	AdoptedAt   *string
	CreatedAt   string
	UpdatedAt   string
}

type CreatePuppyInputPort interface {
	Create(ctx context.Context, cmd CreatePuppyCommand) (PuppyDetails, error)
}

type GetPuppyQuery struct {
	PuppyID string
}

type GetPuppyInputPort interface {
	Get(ctx context.Context, query GetPuppyQuery) (PuppyDetails, error)
}

type ListMyPuppiesQuery struct {
	OwnerID string
	Page    common.PageRequest
}

type ListMyPuppiesInputPort interface {
	ListMine(ctx context.Context, query ListMyPuppiesQuery) (common.Page[PuppyDetails], error)
}

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

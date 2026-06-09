package puppy

import (
	"time"

	"adotapet/internal/domain/common"
)

type Species string
type Size string
type Sex string
type Status string

const (
	SpeciesDog   Species = "DOG"
	SpeciesCat   Species = "CAT"
	SpeciesOther Species = "OTHER"

	SizeSmall  Size = "SMALL"
	SizeMedium Size = "MEDIUM"
	SizeLarge  Size = "LARGE"

	SexMale    Sex = "MALE"
	SexFemale  Sex = "FEMALE"
	SexUnknown Sex = "UNKNOWN"

	StatusAvailable Status = "AVAILABLE"
	StatusAdopted   Status = "ADOPTED"
	StatusPaused    Status = "PAUSED"
	StatusRemoved   Status = "REMOVED"
)

type Puppy struct {
	ID          string
	OwnerID     string
	Name        string
	Breed       *string
	Species     Species
	AgeMonths   int
	Size        Size
	Sex         Sex
	Description string
	Location    common.GeoPoint
	City        string
	State       string
	Status      Status
	AdoptedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

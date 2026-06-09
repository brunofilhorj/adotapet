package puppy

import "time"

type Photo struct {
	ID        string
	PuppyID   string
	URL       string
	SortOrder int
	IsPrimary bool
	CreatedAt time.Time
}

package search

import "time"

type SavedSearch struct {
	ID            string
	UserID        string
	Name          string
	Filters       map[string]any
	NotifyOnMatch bool
	CreatedAt     time.Time
}

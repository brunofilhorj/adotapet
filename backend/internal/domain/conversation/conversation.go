package conversation

import "time"

type Status string

const (
	StatusOpen    Status = "OPEN"
	StatusClosed  Status = "CLOSED"
	StatusBlocked Status = "BLOCKED"
)

type Conversation struct {
	ID        string
	PuppyID   string
	AdopterID string
	DonorID   string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

package user

import "time"

type Role string

const (
	RoleAdopter Role = "ADOPTER"
	RoleDonor   Role = "DONOR"
	RoleShelter Role = "SHELTER"
)

type Status string

const (
	StatusPendingVerification Status = "PENDING_VERIFICATION"
	StatusActive              Status = "ACTIVE"
	StatusSuspended           Status = "SUSPENDED"
	StatusDeleted             Status = "DELETED"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

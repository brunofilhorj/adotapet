package in

import (
	"context"

	"adotapet/internal/domain/common"
)

type MyProfile struct {
	UserID    string
	Email     string
	Role      string
	Status    string
	Name      string
	Phone     *string
	City      string
	State     string
	Location  *common.GeoPoint
	AvatarURL *string
	Bio       *string
}

type GetMyProfileQuery struct {
	UserID string
}

type GetMyProfileInputPort interface {
	Get(ctx context.Context, query GetMyProfileQuery) (MyProfile, error)
}

type UpdateProfileCommand struct {
	UserID    string
	Name      *string
	Phone     *string
	City      *string
	State     *string
	Location  *common.GeoPoint
	AvatarURL *string
	Bio       *string
}

type UpdateProfileInputPort interface {
	Update(ctx context.Context, cmd UpdateProfileCommand) (MyProfile, error)
}

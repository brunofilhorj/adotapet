package user

import "adotapet/internal/domain/common"

type Profile struct {
	UserID    string
	Name      string
	Phone     *string
	City      string
	State     string
	Location  *common.GeoPoint
	AvatarURL *string
	Bio       *string
}

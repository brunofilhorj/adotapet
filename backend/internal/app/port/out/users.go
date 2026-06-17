package out

import (
	"context"

	"adotapet/internal/domain/user"
)

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*user.Profile, error)
	Update(ctx context.Context, profile user.Profile) (user.Profile, error)
}

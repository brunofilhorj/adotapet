package out

import (
	"context"

	"adotapet/internal/domain/common"
	"adotapet/internal/domain/puppy"
	"adotapet/internal/domain/user"
)

type UserRepository interface {
	Save(ctx context.Context, user user.User) (user.User, error)
	SaveWithProfile(ctx context.Context, user user.User, profile user.Profile) (user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Activate(ctx context.Context, id string) (user.User, error)
}

type VerificationCodeRepository interface {
	Save(ctx context.Context, code user.AccountVerificationCode) (user.AccountVerificationCode, error)
	FindPending(ctx context.Context, userID string, channel user.VerificationChannel, destination string, codeHash string) (*user.AccountVerificationCode, error)
	Consume(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token user.RefreshToken) (user.RefreshToken, error)
	FindByHash(ctx context.Context, tokenHash string) (*user.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
}

type PuppyRepository interface {
	Save(ctx context.Context, puppy puppy.Puppy) (puppy.Puppy, error)
	FindByID(ctx context.Context, id string) (*puppy.Puppy, error)
	FindByOwnerID(ctx context.Context, ownerID string, page common.PageRequest) (common.Page[puppy.Puppy], error)
}

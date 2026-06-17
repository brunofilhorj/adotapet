package in

import "context"

type RegisterUserCommand struct {
	Email               string
	Password            string
	Role                string
	Name                string
	City                string
	State               string
	Phone               string
	VerificationChannel string
}

type RegisteredUser struct {
	UserID              string
	Status              string
	VerificationChannel string
	VerificationTarget  string
}

type RegisterUserInputPort interface {
	Register(ctx context.Context, cmd RegisterUserCommand) (RegisteredUser, error)
}

type LoginCommand struct {
	Email    string
	Password string
}

type AuthTokens struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int
	RefreshExpiresIn int
}

type LoginInputPort interface {
	Login(ctx context.Context, cmd LoginCommand) (AuthTokens, error)
}

type RefreshTokenCommand struct {
	RefreshToken string
}

type RefreshTokenInputPort interface {
	Refresh(ctx context.Context, cmd RefreshTokenCommand) (AuthTokens, error)
}

type VerifyAccountCommand struct {
	Email       string
	Channel     string
	Destination string
	Code        string
}

type VerifiedAccount struct {
	UserID string
	Status string
}

type VerifyAccountInputPort interface {
	Verify(ctx context.Context, cmd VerifyAccountCommand) (VerifiedAccount, error)
}

type ResendVerificationCommand struct {
	Email       string
	Channel     string
	Destination string
}

type ResendVerificationResult struct {
	UserID  string
	Channel string
	Target  string
}

type ResendVerificationInputPort interface {
	Resend(ctx context.Context, cmd ResendVerificationCommand) (ResendVerificationResult, error)
}

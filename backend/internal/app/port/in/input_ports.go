package in

import (
	"context"

	"adotapet/internal/domain/common"
)

type RegisterUserCommand struct {
	Email    string
	Password string
	Role     string
	Name     string
	City     string
	State    string
}

type RegisteredUser struct {
	UserID string
	Status string
}

type RegisterUserInputPort interface {
	Register(ctx context.Context, cmd RegisterUserCommand) (RegisteredUser, error)
}

type LoginCommand struct {
	Email    string
	Password string
}

type AuthTokens struct {
	AccessToken string
	ExpiresIn   int
}

type LoginInputPort interface {
	Login(ctx context.Context, cmd LoginCommand) (AuthTokens, error)
}

type PuppySearchQuery struct {
	Latitude     float64
	Longitude    float64
	RadiusKM     float64
	Species      *string
	Breed        *string
	AgeMinMonths *int
	AgeMaxMonths *int
	Size         *string
	Sex          *string
	Page         common.PageRequest
}

type PuppySummary struct {
	ID              string
	Name            string
	Breed           *string
	Species         string
	AgeMonths       int
	Size            string
	Sex             string
	Status          string
	DistanceKM      float64
	City            string
	State           string
	PrimaryPhotoURL *string
}

type SearchPuppiesInputPort interface {
	Search(ctx context.Context, query PuppySearchQuery) (common.Page[PuppySummary], error)
}

type SendMessageCommand struct {
	ConversationID string
	SenderID       string
	Content        string
}

type SentMessage struct {
	MessageID string
	SentAt    string
}

type SendMessageInputPort interface {
	Send(ctx context.Context, cmd SendMessageCommand) (SentMessage, error)
}

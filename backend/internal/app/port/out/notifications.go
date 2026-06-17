package out

import (
	"context"
	"time"

	"adotapet/internal/domain/user"
)

type NotificationPort interface {
	SendAccountVerification(ctx context.Context, target NotificationTarget, code string) error
	SendNewMessage(ctx context.Context, target NotificationTarget, message MessageNotification) error
	SendSavedSearchMatch(ctx context.Context, target NotificationTarget, match SavedSearchMatch) error
}

type VerificationSender interface {
	SendVerificationCode(ctx context.Context, message VerificationMessage) error
}

type VerificationMessage struct {
	UserID      string
	Channel     user.VerificationChannel
	Destination string
	Code        string
	ExpiresAt   time.Time
}

type NotificationTarget struct {
	UserID string
	Email  string
}

type MessageNotification struct {
	ConversationID string
	MessageID      string
}

type SavedSearchMatch struct {
	SavedSearchID string
	PuppyID       string
}

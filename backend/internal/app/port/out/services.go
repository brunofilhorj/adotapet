package out

import (
	"context"
	"time"

	inport "adotapet/internal/app/port/in"
	"adotapet/internal/domain/common"
)

type SearchRepository interface {
	SearchNearby(ctx context.Context, query inport.PuppySearchQuery) (common.Page[inport.PuppySummary], error)
}

type MediaStoragePort interface {
	CreateUploadURL(ctx context.Context, cmd CreateUploadURLCommand) (PresignedUpload, error)
	DeleteObject(ctx context.Context, objectKey string) error
}

type CreateUploadURLCommand struct {
	FileName    string
	ContentType string
	Purpose     string
}

type PresignedUpload struct {
	ObjectKey string
	UploadURL string
	PublicURL string
	ExpiresIn int
}

type GeoLocationPort interface {
	Geocode(ctx context.Context, address Address) (common.GeoPoint, error)
}

type Address struct {
	City   string
	State  string
	Street string
	Number string
}

type CachePort interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	Evict(ctx context.Context, pattern string) error
}

type NotificationPort interface {
	SendAccountVerification(ctx context.Context, target NotificationTarget, code string) error
	SendNewMessage(ctx context.Context, target NotificationTarget, message MessageNotification) error
	SendSavedSearchMatch(ctx context.Context, target NotificationTarget, match SavedSearchMatch) error
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

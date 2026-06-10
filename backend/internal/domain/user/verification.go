package user

import "time"

type VerificationChannel string

const (
	VerificationChannelEmail    VerificationChannel = "EMAIL"
	VerificationChannelSMS      VerificationChannel = "SMS"
	VerificationChannelWhatsApp VerificationChannel = "WHATSAPP"
	VerificationChannelPush     VerificationChannel = "PUSH"
)

type AccountVerificationCode struct {
	ID          string
	UserID      string
	Channel     VerificationChannel
	Destination string
	CodeHash    string
	ExpiresAt   time.Time
	ConsumedAt  *time.Time
	CreatedAt   time.Time
}

func (c AccountVerificationCode) IsUsable(now time.Time) bool {
	return c.ConsumedAt == nil && now.Before(c.ExpiresAt)
}

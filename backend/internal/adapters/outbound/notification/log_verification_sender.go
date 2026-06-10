package notification

import (
	"context"
	"log/slog"

	outport "adotapet/internal/app/port/out"
)

type LogVerificationSender struct {
	log *slog.Logger
}

func NewLogVerificationSender(log *slog.Logger) LogVerificationSender {
	return LogVerificationSender{log: log}
}

func (s LogVerificationSender) SendVerificationCode(ctx context.Context, message outport.VerificationMessage) error {
	s.log.InfoContext(ctx, "verification code issued",
		slog.String("userId", message.UserID),
		slog.String("channel", string(message.Channel)),
		slog.String("destination", message.Destination),
		slog.String("code", message.Code),
		slog.Time("expiresAt", message.ExpiresAt),
	)
	return nil
}

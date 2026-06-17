package in

import "context"

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

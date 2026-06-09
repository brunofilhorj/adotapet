package conversation

import "time"

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Content        string
	SentAt         time.Time
	ReadAt         *time.Time
}

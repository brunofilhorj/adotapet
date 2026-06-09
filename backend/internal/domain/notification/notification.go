package notification

type Channel string

const (
	ChannelPush  Channel = "PUSH"
	ChannelEmail Channel = "EMAIL"
	ChannelSMS   Channel = "SMS"
)

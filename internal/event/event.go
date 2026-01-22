package event

import "github.com/kkrypt0nn/aegisbot/proto"

type EventType string

const (
	EventMessage EventType = "message"
	EventMember  EventType = "member"
)

type Context struct {
	Type EventType

	GuildID   string
	ChannelID string
	MessageID string
	UserID    string

	Message *proto.Message
	Member  *proto.Member
}

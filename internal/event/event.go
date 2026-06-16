package event

import "github.com/kkrypt0nn/aegisbot/proto"

type EventType string

const (
	EventMessageCreate EventType = "message_create"
	EventMessageUpdate EventType = "message_update"
	EventMemberJoin    EventType = "member_join"
	EventMemberUpdate  EventType = "member_update"
)

type Context struct {
	Type EventType

	Bot bool

	GuildID   string
	ChannelID string
	MessageID string
	UserID    string

	Message *proto.Message
	Member  *proto.Member
}

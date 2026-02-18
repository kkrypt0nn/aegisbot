package event

import "github.com/kkrypt0nn/aegisbot/proto"

type EventType string

const (
	EventMessage      EventType = "message"
	EventMemberJoin   EventType = "member_join"
	EventMemberUpdate EventType = "member_update"
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

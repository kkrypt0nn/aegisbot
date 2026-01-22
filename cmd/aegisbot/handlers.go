package main

import (
	"github.com/disgoorg/disgo/events"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/proto"
)

func (b *Bot) handleMessage(e *events.MessageCreate) {
	member := &proto.Member{
		Username: e.Message.Author.Username,
		Bot:      e.Message.Author.Bot,
	}
	ctx := &event.Context{
		Type:      event.EventMessage,
		GuildID:   e.GuildID.String(),
		ChannelID: e.ChannelID.String(),
		MessageID: e.MessageID.String(),
		UserID:    e.Message.Author.ID.String(),
		Message: &proto.Message{
			Content: e.Message.Content,
			Author:  member,
		},
		Member: member,
	}
	b.ProcessRules(ctx)
}

func (b *Bot) handleMemberUpdate(e *events.GuildMemberUpdate) {
	ctx := &event.Context{
		Type:    event.EventMember,
		GuildID: e.GuildID.String(),
		UserID:  e.Member.User.ID.String(),
		Member: &proto.Member{
			Username: e.Member.User.Username,
			Bot:      e.Member.User.Bot,
		},
	}
	b.ProcessRules(ctx)
}

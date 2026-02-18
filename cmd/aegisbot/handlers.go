package main

import (
	"github.com/disgoorg/disgo/events"
	"github.com/kkrypt0nn/aegisbot/internal/commands"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (b *Bot) handleCommand(event *events.ApplicationCommandInteractionCreate) {
	commands.CommandsList[event.SlashCommandInteractionData().CommandName()].Handle(event, b.RulesByName)
}

func (b *Bot) handleMessage(e *events.MessageCreate) {
	member := &proto.Member{
		Username:  e.Message.Author.Username,
		Bot:       e.Message.Author.Bot,
		CreatedAt: timestamppb.New(e.Message.Author.ID.Time()),
	}

	message := &proto.Message{
		Content:  e.Message.Content,
		Author:   member,
		Mentions: make([]string, len(e.Message.Mentions)),
	}
	for i, user := range e.Message.Mentions {
		message.Mentions[i] = user.Username
	}

	discordChannel, exists := e.Channel()
	if exists {
		message.Channel = &proto.Channel{
			Name: discordChannel.Name(),
			Type: proto.ChannelType(discordChannel.Type() + 1),
		}
	}

	ctx := &event.Context{
		Type:      event.EventMessage,
		ChannelID: e.ChannelID.String(),
		MessageID: e.MessageID.String(),
		UserID:    e.Message.Author.ID.String(),
		Message:   message,
		Member:    member,
	}

	if e.GuildID != nil {
		ctx.GuildID = e.GuildID.String()
	}

	b.ProcessRules(ctx)
}

func (b *Bot) handleMemberUpdate(e *events.GuildMemberUpdate) {
	ctx := &event.Context{
		Type:    event.EventMemberUpdate,
		GuildID: e.GuildID.String(),
		UserID:  e.Member.User.ID.String(),
		Member: &proto.Member{
			Username:  e.Member.User.Username,
			Bot:       e.Member.User.Bot,
			CreatedAt: timestamppb.New(e.Member.User.ID.Time()),
		},
	}

	b.ProcessRules(ctx)
}

func (b *Bot) handleMemberJoin(e *events.GuildMemberJoin) {
	ctx := &event.Context{
		Type:    event.EventMemberUpdate,
		GuildID: e.GuildID.String(),
		UserID:  e.Member.User.ID.String(),
		Member: &proto.Member{
			Username:  e.Member.User.Username,
			Bot:       e.Member.User.Bot,
			CreatedAt: timestamppb.New(e.Member.User.ID.Time()),
		},
	}

	b.ProcessRules(ctx)
}

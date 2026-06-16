package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/kkrypt0nn/aegisbot/internal/commands"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (b *Bot) handleCommand(event *events.ApplicationCommandInteractionCreate) {
	commands.CommandsList[event.SlashCommandInteractionData().CommandName()].Handle(event, b.RulesByName)
}

func (b *Bot) handleMessage(e *events.MessageCreate) {
	b.processMessageEvent(
		event.EventMessageCreate,
		e.GuildID,
		e.ChannelID,
		e.MessageID,
		e.Message,
		e.Channel,
	)
}

func (b *Bot) handleMessageUpdate(e *events.MessageUpdate) {
	b.processMessageEvent(
		event.EventMessageUpdate,
		e.GuildID,
		e.ChannelID,
		e.MessageID,
		e.Message,
		e.Channel,
	)
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
			JoinedAt:  timestamppb.New(*e.Member.JoinedAt),
		},
	}

	b.ProcessRules(ctx)
}

func (b *Bot) handleMemberJoin(e *events.GuildMemberJoin) {
	ctx := &event.Context{
		Type:    event.EventMemberJoin,
		Bot:     e.Member.User.Bot,
		GuildID: e.GuildID.String(),
		UserID:  e.Member.User.ID.String(),
		Member: &proto.Member{
			Username:  e.Member.User.Username,
			Bot:       e.Member.User.Bot,
			CreatedAt: timestamppb.New(e.Member.User.ID.Time()),
			JoinedAt:  timestamppb.New(*e.Member.JoinedAt),
		},
	}

	b.ProcessRules(ctx)
}

func (b *Bot) processMessageEvent(
	eventType event.EventType,
	guildID *snowflake.ID,
	channelID snowflake.ID,
	messageID snowflake.ID,
	msg discord.Message,
	getChannel func() (discord.GuildMessageChannel, bool),
) {
	member := &proto.Member{
		Username:  msg.Author.Username,
		Bot:       msg.Author.Bot,
		CreatedAt: timestamppb.New(msg.Author.ID.Time()),
	}

	if msg.Member != nil && msg.Member.JoinedAt != nil {
		member.JoinedAt = timestamppb.New(*msg.Member.JoinedAt)
	}

	message := &proto.Message{
		Content:     msg.Content,
		Author:      member,
		Mentions:    make([]string, len(msg.Mentions)),
		Attachments: make([]*proto.Attachment, len(msg.Attachments)),
	}
	for i, user := range msg.Mentions {
		message.Mentions[i] = user.Username
	}
	for i, attachment := range msg.Attachments {
		message.Attachments[i] = &proto.Attachment{
			Filename: attachment.Filename,
			Url:      attachment.URL,
		}
	}

	if discordChannel, ok := getChannel(); ok {
		message.Channel = &proto.Channel{
			Name: discordChannel.Name(),
			Type: proto.ChannelType(discordChannel.Type() + 1),
		}
	}

	ctx := &event.Context{
		Type:      eventType,
		Bot:       msg.Author.Bot,
		ChannelID: channelID.String(),
		MessageID: messageID.String(),
		UserID:    msg.Author.ID.String(),
		Message:   message,
	}

	if guildID != nil {
		ctx.GuildID = guildID.String()
	}

	b.ProcessRules(ctx)
}

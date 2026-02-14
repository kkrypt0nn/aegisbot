package actions

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/kkrypt0nn/aegisbot/internal/log"
	"github.com/kkrypt0nn/aegisbot/internal/template"
)

type Input struct {
	RuleName string

	GuildID         string
	ChannelID       string
	MessageID       string
	UserID          string
	Duration        time.Duration
	Reason          string
	MessageTemplate string

	Variables map[string]any
}

func Execute(action string, client rest.Rest, input *Input) {
	switch action {
	case "alert":
		alert(client,
			snowflake.MustParse(input.ChannelID),
			template.Render(
				input.MessageTemplate,
				input.Variables,
				template.DefaultAlert,
			),
		)

	case "ban":
		ban(client,
			parseSnowflakePtr(input.GuildID),
			snowflake.MustParse(input.UserID),
			template.Render(
				input.Reason,
				input.Variables,
				template.DefaultBan,
			),
		)

	case "delete":
		deleteMessage(
			client,
			snowflake.MustParse(input.ChannelID),
			snowflake.MustParse(input.MessageID),
		)

	case "kick":
		kick(
			client,
			parseSnowflakePtr(input.GuildID),
			snowflake.MustParse(input.UserID),
			template.Render(
				input.Reason,
				input.Variables,
				template.DefaultKick,
			),
		)

	case "timeout":
		duration := input.Duration
		if duration == 0 {
			duration = 10 * time.Minute
		}
		timeout(
			client,
			parseSnowflakePtr(input.GuildID),
			snowflake.MustParse(input.UserID),
			duration,
		)

	default:
		log.Warn(fmt.Sprintf("unknown action: %s", action))
	}
}

func alert(client rest.Rest, channelID snowflake.ID, message string) {
	_, err := client.CreateMessage(channelID, discord.NewMessageCreate().WithContent(message))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send alert message: %v", err))
	}
}

func ban(client rest.Rest, guildID *snowflake.ID, userID snowflake.ID, reason string) {
	if guildID == nil {
		return
	}

	err := client.AddBan(*guildID, userID, 0, rest.WithReason(reason))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to ban user: %v", err))
	}
}

func deleteMessage(client rest.Rest, channelID, messageID snowflake.ID) {
	err := client.DeleteMessage(channelID, messageID)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to delete message: %v", err))
	}
}

func kick(client rest.Rest, guildID *snowflake.ID, userID snowflake.ID, reason string) {
	if guildID == nil {
		return
	}

	err := client.RemoveMember(*guildID, userID, rest.WithReason(reason))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to kick user: %v", err))
	}
}

func timeout(client rest.Rest, guildID *snowflake.ID, userID snowflake.ID, duration time.Duration) {
	if guildID == nil {
		return
	}

	until := time.Now().Add(duration)
	_, err := client.UpdateMember(*guildID, userID, discord.MemberUpdate{
		CommunicationDisabledUntil: omit.New(&until),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to timeout user: %v", err))
	}
}

func parseSnowflakePtr(id string) *snowflake.ID {
	if id == "" {
		return nil
	}
	snow := snowflake.MustParse(id)
	return &snow
}

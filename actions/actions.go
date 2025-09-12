package actions

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/kkrypt0nn/aegisbot/log"
)

type Context struct {
	GuildID   string
	ChannelID string
	UserID    string
	MessageID string
	Duration  time.Duration
}

func Execute(action string, client rest.Rest, ctx *Context) {
	switch action {
	case "alert":
		alert(client, snowflake.MustParse(ctx.ChannelID), snowflake.MustParse(ctx.UserID))

	case "ban":
		guildID := parseSnowflakePtr(ctx.GuildID)
		ban(client, guildID, snowflake.MustParse(ctx.UserID))

	case "delete":
		deleteMessage(client,
			snowflake.MustParse(ctx.ChannelID),
			snowflake.MustParse(ctx.MessageID))

	case "kick":
		guildID := parseSnowflakePtr(ctx.GuildID)
		kick(client, guildID, snowflake.MustParse(ctx.UserID))

	case "timeout":
		guildID := parseSnowflakePtr(ctx.GuildID)
		duration := ctx.Duration
		if duration == 0 {
			duration = 10 * time.Minute
		}
		timeout(client, guildID, snowflake.MustParse(ctx.UserID), duration)

	default:
		log.Warn(fmt.Sprintf("unknown action: %s", action))
	}
}

func alert(client rest.Rest, channelID, userID snowflake.ID) {
	// Yes, this is hard-coded and for now and I'm fine with it...
	content := fmt.Sprintf("<@&861907005386915870> ⚠️ Rule has been matched and triggered by <@%s>", userID)
	_, err := client.CreateMessage(channelID, discord.NewMessageCreateBuilder().SetContent(content).Build())
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send alert message: %v", err))
	}
}

func ban(client rest.Rest, guildID *snowflake.ID, userID snowflake.ID) {
	if guildID == nil {
		return
	}

	err := client.AddBan(*guildID, userID, 0, rest.WithReason("AutoMod: Rule matched"))
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

func kick(client rest.Rest, guildID *snowflake.ID, userID snowflake.ID) {
	if guildID == nil {
		return
	}

	err := client.RemoveMember(*guildID, userID, rest.WithReason("AutoMod: Rule matched"))
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
		CommunicationDisabledUntil: json.NewNullablePtr(until),
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

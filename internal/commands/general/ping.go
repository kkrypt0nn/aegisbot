package cmdgeneral

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/kkrypt0nn/aegisbot/internal/commands/cmdtypes"
	"github.com/kkrypt0nn/aegisbot/internal/log"
	"github.com/kkrypt0nn/aegisbot/internal/rules"
)

type Ping struct {
	Name        string
	Description string
	Options     []discord.ApplicationCommandOption
	Category    cmdtypes.Category
}

func (c Ping) CommandCreateData() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        c.Name,
		Description: c.Description,
		Options:     c.Options,
	}
}

func (c Ping) Handle(event *events.ApplicationCommandInteractionCreate, _ map[string]*rules.SimplifiedRule) {
	err := event.CreateMessage(discord.NewMessageCreate().
		WithContent("Hi there, I'm ready to watch people trigger your rules!").
		WithEphemeral(true),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send response: %v", err))
	}
}

func (c Ping) Usage() string {
	options := ""
	if c.Options != nil {
		for _, option := range c.Options {
			options += fmt.Sprintf("<%s> ", option.OptionName())
		}
	}
	return "/ping " + options
}

package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/kkrypt0nn/aegisbot/internal/commands/cmdtypes"
	cmdgeneral "github.com/kkrypt0nn/aegisbot/internal/commands/general"
	cmdrules "github.com/kkrypt0nn/aegisbot/internal/commands/rules"
)

var CommandsList = map[string]cmdtypes.Command{
	// General
	"ping": cmdgeneral.Ping{
		Name:        "ping",
		Description: "Check if the bot is alive",
		Options:     []discord.ApplicationCommandOption{},
		Category:    cmdtypes.CategoryGeneral,
	},

	// Rules
	"get": cmdrules.Get{
		Name:        "get",
		Description: "Get the YAML source of a given rule name",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "name",
				Description: "The unique rule name",
				Required:    true,
			},
		},
		Category: cmdtypes.CategoryRules,
	},
}

func PrepareCommandCreateData() []discord.ApplicationCommandCreate {
	slashCommands := make([]discord.ApplicationCommandCreate, 0, len(CommandsList))
	for _, command := range CommandsList {
		slashCommands = append(slashCommands, command.CommandCreateData())
	}
	return slashCommands
}

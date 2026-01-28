package cmdtypes

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/kkrypt0nn/aegisbot/internal/rules"
)

type Category int

const (
	CategoryGeneral Category = iota
	CategoryRules
)

type Command interface {
	Handle(*events.ApplicationCommandInteractionCreate, map[string]*rules.SimplifiedRule)
	CommandCreateData() discord.SlashCommandCreate
	Usage() string
}

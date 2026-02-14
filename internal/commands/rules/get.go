package cmdrules

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/kkrypt0nn/aegisbot/internal/commands/cmdtypes"
	"github.com/kkrypt0nn/aegisbot/internal/log"
	"github.com/kkrypt0nn/aegisbot/internal/rules"
)

type Get struct {
	Name        string
	Description string
	Options     []discord.ApplicationCommandOption
	Category    cmdtypes.Category
}

func (c Get) CommandCreateData() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        c.Name,
		Description: c.Description,
		Options:     c.Options,
	}
}

func (c Get) Handle(event *events.ApplicationCommandInteractionCreate, rulesByName map[string]*rules.SimplifiedRule) {
	ruleName := event.SlashCommandInteractionData().String("name")
	rule, ok := rulesByName[ruleName]
	if !ok {
		err := event.CreateMessage(discord.NewMessageCreate().
			WithContent(fmt.Sprintf("❌ Rule `%s` does not exist!", ruleName)).
			WithEphemeral(true),
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to send response: %v", err))
		}
		return
	}

	description := "### Meta\n"
	description += fmt.Sprintf("- **Name**: %s\n", rule.Name)
	description += fmt.Sprintf("- **Event**: %s\n", rule.Event)
	description += fmt.Sprintf("- **Ignore bots**: %s\n", fmt.Sprintf("%v", rule.IgnoreBots))
	description += fmt.Sprintf("- **Action**: %s\n", string(rule.Action.Type))

	description += "\n### Strings\n"

	if len(rule.Strings) > 0 {
		for k, v := range rule.Strings {
			description += fmt.Sprintf("`%s` -> `%s`\n", k, v.Value)
		}
	}

	description += "\n### Expression\n"
	description += fmt.Sprintf("```yaml\n%s\n```", rule.RawExpression)

	err := event.CreateMessage(discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbedBuilder().
			SetDescription(description).
			Build(),
		).
		WithEphemeral(true),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send response: %v", err))
	}
}

func (c Get) Usage() string {
	options := ""
	if c.Options != nil {
		for _, option := range c.Options {
			options += fmt.Sprintf("<%s> ", option.OptionName())
		}
	}
	return "/get " + options
}

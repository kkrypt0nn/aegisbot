package main

import (
	"fmt"

	"github.com/disgoorg/disgo/bot"
	"github.com/kkrypt0nn/aegisbot/internal/actions"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/internal/log"
	"github.com/kkrypt0nn/aegisbot/internal/rules"
)

type Bot struct {
	Client      *bot.Client
	Config      *Config
	Rules       []*rules.SimplifiedRule
	RulesByName map[string]*rules.SimplifiedRule
}

type Config struct {
	IgnoreBots  bool
	RulesFolder string
}

func (b *Bot) ProcessRules(ctx *event.Context) {
	for _, rule := range b.Rules {
		ok, err := rule.Evaluate(ctx)
		if err != nil || !ok {
			continue
		}

		// Prepare the template variables
		variables := map[string]any{
			"RuleName":  rule.Name,
			"GuildID":   ctx.GuildID,
			"ChannelID": ctx.ChannelID,
			"MessageID": ctx.MessageID,
			"UserID":    ctx.UserID,
		}
		for k, s := range rule.Strings {
			variables[k] = s.Value
		}

		log.Info(fmt.Sprintf("Rule matched: %s", rule.Name))

		actions.Execute(string(rule.Action.Type), b.Client.Rest, &actions.Input{
			RuleName: rule.Name,

			GuildID:         ctx.GuildID,
			ChannelID:       ctx.ChannelID,
			MessageID:       ctx.MessageID,
			UserID:          ctx.UserID,
			Reason:          rule.Action.Reason,
			MessageTemplate: rule.Action.MessageTemplate,

			Variables: variables,
		})
	}
}

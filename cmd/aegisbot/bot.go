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
	Client bot.Client
	Config *Config
	Rules  []*rules.SimplifiedRule
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

		actions.Execute(rule.Action, b.Client.Rest(), &actions.Input{
			RuleName: rule.Name,

			GuildID:   ctx.GuildID,
			ChannelID: ctx.ChannelID,
			MessageID: ctx.MessageID,
			UserID:    ctx.UserID,

			Variables: variables,

			AlertTemplate: rule.AlertTemplate,
			BanTemplate:   rule.BanTemplate,
			KickTemplate:  rule.KickTemplate,
		})
	}
}

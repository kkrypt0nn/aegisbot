package main

import (
	"time"

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

		log.Infof("Rule matched: %s", rule.Name)

		for _, action := range rule.Actions {
			var timeoutDuration time.Duration
			if action.Duration != "" {
				duration, err := time.ParseDuration(action.Duration)
				if err != nil {
					log.Warnf(
						"Invalid duration (%s) for rule %s: %v",
						action.Duration,
						rule.Name,
						err,
					)
				} else {
					timeoutDuration = duration
				}
			}

			actions.Execute(action.Type, b.Client.Rest, &actions.Input{
				RuleName:  rule.Name,
				GuildID:   ctx.GuildID,
				ChannelID: ctx.ChannelID,
				MessageID: ctx.MessageID,
				UserID:    ctx.UserID,

				AlertChannelID:  action.ChannelID,
				AlertMessage:    action.Message,
				BanKickReason:   action.Reason,
				TimeoutDuration: timeoutDuration,

				Variables: variables,
			})
		}
	}
}

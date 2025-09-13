package main

import "github.com/disgoorg/disgo/events"

func (b *Bot) returnIfBot(bot bool) bool {
	return b.Config.IgnoreBots && bot
}

func (b *Bot) handleMessage(e *events.MessageCreate) {
	if b.returnIfBot(e.Message.Author.Bot) {
		return
	}

	for _, rule := range b.Rules {
		rule.EvaluateMessage(e)
	}
}

func (b *Bot) handleMemberUpdate(e *events.GuildMemberUpdate) {
	if b.returnIfBot(e.Member.User.Bot) {
		return
	}

	for _, rule := range b.Rules {
		rule.EvaluateMember(e)
	}
}

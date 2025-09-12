package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/kkrypt0nn/aegisbot/core"
	"github.com/kkrypt0nn/aegisbot/log"
	"github.com/kkrypt0nn/aegisbot/rules"
)

type Bot struct {
	Client bot.Client
	Rules  []*rules.SimplifiedRule
}

func (b *Bot) handleMessage(e *events.MessageCreate) {
	if e.Message.Author.Bot {
		return
	}

	for _, rule := range b.Rules {
		rule.EvaluateMessage(e)
	}
}

func (b *Bot) handleMemberUpdate(e *events.GuildMemberUpdate) {
	for _, rule := range b.Rules {
		rule.EvaluateMember(e)
	}
}

func main() {
	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")

	yarabotBot := &Bot{}

	client, err := disgo.New(token, bot.WithGatewayConfigOpts(
		gateway.WithIntents(
			gateway.IntentGuildMessages,
			gateway.IntentMessageContent,
			gateway.IntentGuildMembers,
		)),
		bot.WithEventListenerFunc(yarabotBot.handleMessage),
		bot.WithEventListenerFunc(yarabotBot.handleMemberUpdate),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed creating Discord session: %s", err))
		return
	}
	yarabotBot.Client = client

	rulesFolder := "_rules/"

	loadedRules, err := rules.Load(rulesFolder)
	if err != nil {
		log.Error(fmt.Sprintf("Failed loading rules: %s", err))
		return
	}

	yarabotBot.Rules = loadedRules

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(fmt.Sprintf("Failed creating file watcher: %s", err))
			return
		}
		defer func() {
			_ = watcher.Close()
		}()

		err = watcher.Add(rulesFolder)
		if err != nil {
			log.Error(fmt.Sprintf("Failed watching rules folder: %s", err))
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					log.Info("Changes detected in rules folder, performing hot-reload...")
					updatedRules, err := rules.Load(rulesFolder)
					if err != nil {
						log.Error(fmt.Sprintf("Failed to reload rules: %s", err))
					} else {
						yarabotBot.Rules = updatedRules
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(fmt.Sprintf("Watcher error: %s", err))
			}
		}
	}()

	defer func() {
		client.Close(context.TODO())
	}()

	err = client.OpenGateway(context.TODO())
	if err != nil {
		log.Error(fmt.Sprintf("Failed opening connection: %s", err))
		return
	}

	log.Success(fmt.Sprintf("%s (v%s) has successfully started", core.Name, core.Version))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}

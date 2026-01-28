package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/kkrypt0nn/aegisbot/internal/buildinfo"
	"github.com/kkrypt0nn/aegisbot/internal/commands"
	"github.com/kkrypt0nn/aegisbot/internal/log"
	"github.com/kkrypt0nn/aegisbot/internal/rules"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")

	config := &Config{}
	aegisbot := &Bot{
		Config: config,
	}

	ignoreBots := os.Getenv("IGNORE_BOTS")
	if strings.ToLower(ignoreBots) == "true" {
		config.IgnoreBots = true
	}

	client, err := disgo.New(token, bot.WithGatewayConfigOpts(
		gateway.WithIntents(
			gateway.IntentGuildMessages,
			gateway.IntentMessageContent,
			gateway.IntentGuildMembers,
		)),
		bot.WithEventListenerFunc(aegisbot.handleCommand),
		bot.WithEventListenerFunc(aegisbot.handleMessage),
		bot.WithEventListenerFunc(aegisbot.handleMemberUpdate),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed creating Discord session: %s", err))
		return
	}
	aegisbot.Client = client

	// TODO: Make this check for commands before, like I did in my personal bot to prevent rate-limits
	if _, err := client.Rest().SetGlobalCommands(client.ApplicationID(), commands.PrepareCommandCreateData()); err != nil {
		log.Error(fmt.Sprintf("Failed to register commands: %v", err))
	}

	rulesFolder := "_rules/"
	config.RulesFolder = rulesFolder

	loadedRules, rulesByName, err := rules.Load(rulesFolder)
	if err != nil {
		log.Error(fmt.Sprintf("Failed loading rules: %s", err))
		return
	}

	aegisbot.Rules = loadedRules
	aegisbot.RulesByName = rulesByName

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
					updatedRules, updatedRulesByName, err := rules.Load(rulesFolder)
					if err != nil {
						log.Error(fmt.Sprintf("Failed to reload rules: %s", err))
					} else {
						aegisbot.Rules = updatedRules
						aegisbot.RulesByName = updatedRulesByName
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

	log.Success(fmt.Sprintf("%s (v%s) has successfully started", buildinfo.Name, buildinfo.Version))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}

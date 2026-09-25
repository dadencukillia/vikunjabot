package main

import (
	_ "github.com/joho/godotenv/autoload"
	"context"
	"fmt"
	"log"
	"vikunjabot/internal"
	"vikunjabot/internal/bot"
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/diffsummary"
	"vikunjabot/internal/texts"
	"vikunjabot/internal/webhook"
)

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	localePack, err := texts.LoadLocalePack(config.Language)
	if err != nil {
		log.Panic(err)
	}

	sumGenerator := diffsummary.NewSummariesGenerator(localePack)
	tgbot := bot.NewBot(config.TelegramBotToken, config.TelegramChatId)

	evLogSq := diffslog.NewDiffsLog()

	server := webhook.NewWebhookServer(config.ServerHost, config.VikunjaWebhookSecret)
	server.Run(context.Background(), func(message webhook.WebhookMessage) error {
		evLogSq.AddEvent(&message)
		evLogSq.Squash()
		tree := diffstree.LogFlowToDiffsTree(evLogSq.GenerateLogFlow())

		for _, text := range sumGenerator.GenerateHTMLSummaries(&tree) {
			resp, err := tgbot.SendTextMessage(text, bot.HTMLParseMode, false, false)
			if err != nil {
				log.Panic(err)
			}

			fmt.Println(resp)
		}

		return nil
	})
}

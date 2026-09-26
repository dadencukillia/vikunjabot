package main

import (
	"context"
	"log"
	"vikunjabot/internal"
	"vikunjabot/internal/bot"
	"vikunjabot/internal/texts"
	"vikunjabot/internal/webhook"

	_ "github.com/joho/godotenv/autoload"
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

	tgbot := bot.NewBot(config.TelegramBotToken, config.TelegramChatId)
	getMeResp, err := tgbot.GetMe()
	if err != nil {
		log.Panic(err)
	}

	if !getMeResp.Ok {
		log.Panicf("please check your telegram token: %s", getMeResp.ErrorDescription)
	}
	log.Printf("Telegram bot logged in as @%s\n", getMeResp.Result.Username)

	pipe := internal.NewPipe(config, localePack, tgbot)

	server := webhook.NewWebhookServer(config.ServerHost, config.VikunjaWebhookSecret)
	server.Run(context.Background(), func(message webhook.WebhookMessage) error {
		pipe.PushEvent(message)

		return nil
	})
}

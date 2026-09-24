package main

import (
	"context"
	"log"
	"vikunjabot/internal"
	"vikunjabot/internal/webhook"
)

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	server := webhook.NewWebhookServer(config.ServerHost, config.ServerHost)
	server.Run(context.Background(), func(message webhook.WebhookMessage) error {
		return nil
	})
}

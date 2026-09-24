package main

import (
	"context"
	"fmt"
	"log"
	"vikunjabot/internal"
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/webhook"

	"github.com/goccy/go-json"
)

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	evLog := diffslog.NewDiffsLog()
	evLogSq := diffslog.NewDiffsLog()

	server := webhook.NewWebhookServer(config.ServerHost, config.VikunjaWebhookSecret)
	server.Run(context.Background(), func(message webhook.WebhookMessage) error {
		evLog.AddEvent(&message)
		evLogSq.AddEvent(&message)
		evLogSq.Squash()

		b, _ := json.Marshal(evLog.GenerateLogFlow())
		fmt.Println("Unsquashed:", string(b))

		bsq, _ := json.Marshal(evLogSq.GenerateLogFlow())
		fmt.Println("Squashed:", string(bsq))

		return nil
	})
}

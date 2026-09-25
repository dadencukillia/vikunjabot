package main

import (
	"context"
	"fmt"
	"log"
	"vikunjabot/internal"
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/diffsummary"
	"vikunjabot/internal/webhook"
)

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	localePack, err := diffsummary.LoadLocalePack(config.Language)
	if err != nil {
		log.Panic(err)
	}

	sumGenerator := diffsummary.NewSummariesGenerator(localePack)

	evLogSq := diffslog.NewDiffsLog()

	server := webhook.NewWebhookServer(config.ServerHost, config.VikunjaWebhookSecret)
	server.Run(context.Background(), func(message webhook.WebhookMessage) error {
		evLogSq.AddEvent(&message)
		evLogSq.Squash()
		tree := diffstree.LogFlowToDiffsTree(evLogSq.GenerateLogFlow())
		fmt.Println(sumGenerator.GenerateHTMLSummaries(&tree))

		return nil
	})
}

package internal

import (
	"log"
	"sync"
	"time"
	"vikunjabot/internal/bot"
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/diffsummary"
	"vikunjabot/internal/texts"
	"vikunjabot/internal/utils"
	"vikunjabot/internal/webhook"
)

type Pipe struct {
	config *Config
	lang *texts.LocalePack
	tgbot *bot.Bot

	mutex *sync.Mutex
	stacks map[int64]*utils.StackDebouncer[webhook.WebhookMessage]
}

func NewPipe(config *Config, lang *texts.LocalePack, tgbot *bot.Bot) Pipe {
	return Pipe{
		config: config,
		lang: lang,
		tgbot: tgbot,

		mutex: &sync.Mutex{},
		stacks: map[int64]*utils.StackDebouncer[webhook.WebhookMessage]{},
	}
}

func (a *Pipe) PushEvent(ev webhook.WebhookMessage) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	projectId := ev.Data.GetProjectID()
	if deb, ok := a.stacks[projectId]; ok {
		deb.Push(ev)
	} else {
		deb := utils.NewStackDebouncer(
			time.Duration(a.config.DebounceSeconds) * time.Second,
			a.RunPipe,
		)
		deb.Push(ev)
		a.stacks[projectId] = deb
	}

	log.Println("Collected!")
}

func (a Pipe) RunPipe(events []webhook.WebhookMessage) {
	logs := diffslog.NewDiffsLog()
	for _, event := range events {
		logs.AddEvent(&event)
	}

	logs.Squash()
	flow := logs.GenerateLogFlow()
	tree := diffstree.LogFlowToDiffsTree(flow)

	textSummGenerator := diffsummary.NewSummariesGenerator(a.lang, a.config.VikunjaHost, a.config.TimeZone)
	summs := textSummGenerator.GenerateHTMLSummaries(&tree)

	for _, text := range summs {
		resp, err := a.tgbot.SendTextMessage(text, bot.HTMLParseMode, false, true)
		if err != nil {
			log.Printf("Something went wrong with sending: %v\n", err)
			continue
		}

		if !resp.Ok {
			log.Printf("Something went wrong with sending: %s\n", resp.ErrorDescription)
			continue
		}

		log.Println("Sent!")
	}
}

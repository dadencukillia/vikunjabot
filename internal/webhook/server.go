package webhook

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type WebhookServer struct {
	serverHost string
	webhookSecret string
}

func NewWebhookServer(host string, secret string) *WebhookServer {
	return &WebhookServer{
		serverHost: host,
		webhookSecret: secret,
	}
}

func (a *WebhookServer) Run(
	ctx context.Context,
	handler func(message WebhookMessage) error,
) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("POST /webhook: %s", r.RemoteAddr)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("no body"))
			return
		}

		fmt.Println(string(body))

		if err := a.verifyHeader(r, body); err != nil {
			log.Printf("signature validation error: %v\n", err)
			w.WriteHeader(400)
			w.Write([]byte("invalid signature"))
			return
		}

		var message WebhookMessage
		err = json.Unmarshal(body, &message)
		if err != nil {
			log.Printf("%v: %v", ErrWebInvalidEventRequest, err)
			w.WriteHeader(400)
			w.Write([]byte(ErrWebInvalidEventRequest.Error()))
			return
		}

		filled := strings.Builder{}
		if message.Data.Task != nil { filled.WriteString("task ") }
		if len(message.Data.Tasks) != 0 { filled.WriteString("tasks ") }
		if message.Data.Assignee != nil { filled.WriteString("assignee ") }
		if message.Data.Comment != nil { filled.WriteString("comment ") }
		if message.Data.Attachment != nil { filled.WriteString("attachment ") }
		if message.Data.Relation != nil { filled.WriteString("relation ") }
		if message.Data.Reminder != nil { filled.WriteString("reminder ") }
		if message.Data.Project != nil { filled.WriteString("project ") }
		if len(message.Data.Projects) != 0 { filled.WriteString("projects ") }
		if message.Data.User != nil { filled.WriteString("user ") }
		if message.Data.Team != nil { filled.WriteString("team ") }

		log.Printf("Got event %s from %s: doer %s", message.EventName, message.Data.Doer.Username, filled.String())

		handler(message)
	})

	s := &http.Server{
		Addr: a.serverHost,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	return s.ListenAndServe()
}

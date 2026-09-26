package internal

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN"`
	TelegramChatId int64 `env:"TELEGRAM_CHAT_ID"`
	VikunjaWebhookSecret string `env:"VIKUNJA_WEBHOOK_SECRET"`
	VikunjaHost string `env:"VIKUNJA_HOST" envDefault:"http://localhost:4321"`
	Language string `env:"SUMMARY_LANG" envDefault:"en"`
	ServerHost string `env:"SERVER_HOST" envDefault:"0.0.0.0:8080"`
}

var parsedConfig *Config = nil

func GetConfig() (*Config, error) {
	if parsedConfig != nil {
		return parsedConfig, nil
	}

	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	parsedConfig = config

	return parsedConfig, nil
}

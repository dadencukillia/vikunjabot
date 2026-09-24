package internal

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN"`
	TelegramChatId int64 `env:"TELEGRAM_CHAT_ID"`
	VikunjaWebhookSecret string `env:"VIKUNJA_WEBHOOK_SECRET"`
	Language string `env:"LANG" envDefault:"en"`
	ServerHost string `env:"SERVER_HOST" envDefault:"0.0.0.0:8080"`
}

var parsedConfig *Config = &Config{}

func GetConfig() (*Config, error) {
	err := env.Parse(parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	return parsedConfig, nil
}

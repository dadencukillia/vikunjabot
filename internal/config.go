package internal

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN"`
	TelegramChatId int64 `env:"TELEGRAM_CHAT_ID"`
	VikunjaWebhookSecret string `env:"VIKUNJA_WEBHOOK_SECRET"`
	Language string `env:"LANG"`
}

var parsedConfig *Config = nil

func GetConfig() (*Config, error) {
	if parsedConfig != nil {
		return parsedConfig, nil
	}

	err := env.Parse(parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	return parsedConfig, nil
}

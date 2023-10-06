package main

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	applogger "telegram-bot-template/pkg/application/logger"
	"telegram-bot-template/pkg/infrastructure/transport/telegram"
)

func messageHandler(ctx context.Context, config *Config, logger applogger.Logger) error {
	api, err := tgbotapi.NewBotAPI(config.AccessToken)
	if err != nil {
		return errors.Wrap(err, "failed to init telegram api")
	}

	handler := telegram.NewEchoHandler()
	messageHandlerConfig := telegram.Config{
		Offset:         config.UpdateConfig.Offset,
		Limit:          config.UpdateConfig.Limit,
		Timeout:        config.UpdateConfig.Timeout,
		AllowedUpdates: config.UpdateConfig.AllowedUpdates,
	}
	tgMessageHandler := telegram.NewMessageHandler(api, config.SleepInterval, messageHandlerConfig, handler, logger)
	return tgMessageHandler.Start(ctx)
}

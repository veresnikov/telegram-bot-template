package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewEchoHandler() Handler {
	return &echoHandler{}
}

type echoHandler struct{}

func (e echoHandler) Handle(_ context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	response.ReplyToMessageID = update.Message.MessageID

	return &response, nil
}

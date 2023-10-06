package telegram

import (
	"context"
	stderrors "errors"
	"time"

	applogger "telegram-bot-template/pkg/application/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrMessageHandlerIsStopped = stderrors.New("message handler is stopped")
)

type Config tgbotapi.UpdateConfig

type Handler interface {
	Handle(ctx context.Context, update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

func NewMessageHandler(
	api *tgbotapi.BotAPI,
	sleepInterval time.Duration,
	config Config,
	handler Handler,
	logger applogger.Logger,
) *MessageHandler {
	return &MessageHandler{
		api:           api,
		sleepInterval: sleepInterval,
		config:        config,
		handler:       handler,
		logger:        logger,
	}
}

type MessageHandler struct {
	api           *tgbotapi.BotAPI
	config        Config
	sleepInterval time.Duration

	handler Handler
	logger  applogger.Logger
}

func (handler *MessageHandler) Start(ctx context.Context) error {
	updates := handler.api.GetUpdatesChan(tgbotapi.UpdateConfig(handler.config))
	var isShutdown bool
	for {
		select {
		case <-ctx.Done():
			if !isShutdown {
				handler.api.StopReceivingUpdates()
				isShutdown = true
			}
		case update, ok := <-updates:
			if ok {
				handler.handle(update)
			} else {
				return ErrMessageHandlerIsStopped
			}
			continue
		default:
			<-time.After(handler.sleepInterval)
		}
	}
}

func (handler *MessageHandler) handle(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	start := time.Now()
	response, handleError := handler.handler.Handle(context.Background(), update)
	var sendError error
	if response != nil {
		_, sendError = handler.api.Send(*response)
	}
	handleTime := time.Since(start)
	fields := applogger.Fields{
		"messageID":   update.Message.MessageID,
		"messageText": update.Message.Text,
		"username":    update.Message.From.UserName,
		"time":        handleTime.String(),
	}
	l := handler.logger.WithFields(fields)
	if handleError == nil && sendError == nil {
		l.Info("call finished")
	} else {
		var err error
		if handleError != nil {
			err = appendToError(err, handleError)
		}
		if sendError != nil {
			err = appendToError(err, sendError)
		}
		l.Error(err, "call failed")
	}
}

func appendToError(err, next error) error {
	if err == nil {
		return next
	}
	return stderrors.Join(err, next)
}

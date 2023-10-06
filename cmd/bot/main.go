package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"telegram-bot-template/pkg/infrastructure/logger"
	"telegram-bot-template/pkg/infrastructure/transport/telegram"
)

const (
	applicationID = "telegram-bot-template"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	ctx = listenOSKillSignalsContext(ctx)
	mainLogger := logger.NewLogger(&logger.Config{AppName: applicationID})
	config, err := parseEnv()
	if err != nil {
		mainLogger.FatalError(err, "failed to parse env")
	}
	err = messageHandler(ctx, config, mainLogger)
	if errors.Cause(err) == telegram.ErrMessageHandlerIsStopped {
		mainLogger.Info(err)
	} else {
		mainLogger.FatalError(err)
	}
}

func listenOSKillSignalsContext(ctx context.Context) context.Context {
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(ctx)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ch:
			cancelFunc()
		case <-ctx.Done():
			return
		}
	}()
	return ctx
}

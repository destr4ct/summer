package main

import (
	"context"
	"destr4ct/summer/internal/config"
	"destr4ct/summer/internal/crawler"
	"destr4ct/summer/internal/queue/rmq"
	"destr4ct/summer/internal/storage/postgres"
	"destr4ct/summer/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	cfg := config.Load()
	logger := logging.GetLogger(cfg.Env)

	broker, err := rmq.GetBroker(cfg)
	if err != nil {
		logger.Error("failed to get broker: %v\n", err)
		return
	}

	defer broker.Close()
	logger.Info("initialized the broker")

	storage, err := postgres.GetStorage(&cfg.DBConfig)
	if err != nil {
		logger.Error("failed to get storage: %v\n", err)
		return
	}
	logger.Info("initialized the storage")

	service := crawler.GetService(logger, storage, broker)
	logger.Info("starting the crawler")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := service.Run(ctx, cfg.Delay); err != nil {
			logger.Error("failed to crawl", err)
		}
		exitSignal <- syscall.SIGTERM
	}()

	<-exitSignal
	cancel()
}

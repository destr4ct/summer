package main

import (
	"context"
	"destr4ct/summer/internal/config"
	"destr4ct/summer/internal/queue"
	"destr4ct/summer/internal/queue/rmq"
	"destr4ct/summer/internal/storage/postgres"
	"destr4ct/summer/pkg/logging"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

	ctx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-exitSignal:
			cancel()
		default:
			bCtx, timeout := context.WithTimeout(ctx, cfg.Delay)

			//logger.Info("consumer: getting the messages")
			messages, err := broker.GetMessages(bCtx, queue.SummerKey)
			if err != nil {
				logger.Error("failed to get messages", err)
			}

			for _, m := range messages {
				// TODO: это лишь симуляция для проверки crawler, нужно поменять потом все к чертям
				aID, _ := strconv.Atoi(m.Message)
				_ = storage.AddSummary(ctx, aID, fmt.Sprintf("test summary for %d", aID))
			}

			timeout()
		}
	}
}

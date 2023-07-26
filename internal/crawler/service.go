package crawler

import (
	"context"
	"destr4ct/summer/internal/queue"
	"destr4ct/summer/internal/storage"
	"golang.org/x/exp/slog"
	"time"
)

type Crawler struct {
	logger  *slog.Logger
	storage storage.SummerStorage
	broker  queue.MessageBroker[time.Time]
}

func (cs *Crawler) Run(ctx context.Context, delay time.Duration) error {
	for {
		select {
		case <-ctx.Done():
		
		}
	}

	return nil
}

func GetService(logger *slog.Logger, storage storage.SummerStorage, broker queue.MessageBroker[time.Time]) *Crawler {
	return &Crawler{
		logger:  logger,
		storage: storage,
		broker:  broker,
	}
}

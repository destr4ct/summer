package crawler

import (
	"context"
	"destr4ct/summer/internal/queue"
	"destr4ct/summer/internal/storage"
	"destr4ct/summer/pkg/utils"
	"errors"
	"golang.org/x/exp/slog"
	"strconv"
	"sync"
	"time"
)

type Crawler struct {
	logger     *slog.Logger
	storage    storage.SummerStorage
	broker     queue.MessageBroker[time.Time]
	maxWorkers int
}

func (cs *Crawler) Run(ctx context.Context, delay time.Duration) error {
	// раз в delay
	//	получаем список ресурсов, по которым нужно пройти
	// 	на каждый ресурс натравливаем краулер
	//  получаем статьи и записываем в базу данных, id-шники отправляем брокеру

	ticker := time.Tick(delay)

	for {
		select {

		case <-ticker:
			cs.logger.Info("performing crawler iteration")

			// Делаем 5 попыток получить sources
			sources, err := utils.DoWithAttempts(func() ([]string, error) {
				return cs.storage.GetAllSources(ctx)
			}, 5)
			if err != nil {
				return err
			}

			cs.process(ctx, sources)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *Crawler) process(ctx context.Context, sources []string) {
	var wg = new(sync.WaitGroup)

	taskPool := make(chan string, len(sources)+1)

	wg.Add(cs.maxWorkers)
	for i := 0; i < cs.maxWorkers; i += 1 {
		cs.logger.Info("started the worker", slog.Int("wID", i))

		wCtx := context.WithValue(ctx, "wID", i)
		go cs.worker(wCtx, taskPool, wg)
	}

	for _, source := range sources {
		taskPool <- source
	}
	close(taskPool)

	wg.Wait()
}

func (cs *Crawler) worker(ctx context.Context, pool chan string, wg *sync.WaitGroup) {
	wID := ctx.Value("wID").(int)

	for {
		select {
		case rawSource, ok := <-pool:
			if !ok {
				cs.logger.Info("Done", slog.Int("wID", wID))
				wg.Done()
				return
			}

			source := getSource(rawSource)
			cs.logger.Info("got source", slog.String("source", rawSource), slog.Int("wID", wID))

			if hdl, ok := handlers[source.Kind]; ok {
				articles, err := utils.DoWithAttempts(func() ([]*storage.SArticle, error) {
					return hdl.ParseSource(ctx, source.Link)
				}, 3)

				if err != nil {
					cs.logger.Error("failed to parse source (3 attempts)", err)
					wg.Done()
					return
				}

				// Сохраняем статьи
				for _, article := range articles {
					updatedArticle, err := cs.storage.AddArticle(ctx, article.Source, article.Content, article.Title)

					if err != nil {
						if !errors.Is(err, storage.ErrDuplicate) {
							cs.logger.Error("failed to save article", err)
							continue
						}
					}

					if updatedArticle.HasSummary {
						continue
					}

					err = cs.broker.SendMessage(ctx, queue.SummerKey, queue.Message[any]{
						Message: strconv.Itoa(updatedArticle.ID),
						Other:   time.Now(),
					})
					if err != nil {
						cs.logger.Error("failed to push article id", err)
					}
				}
				cs.logger.Info("done with source", slog.String("source", rawSource), slog.Int("wID", wID))
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

func GetService(logger *slog.Logger, storage storage.SummerStorage, broker queue.MessageBroker[time.Time], mw int) *Crawler {
	return &Crawler{
		logger:     logger,
		storage:    storage,
		broker:     broker,
		maxWorkers: mw,
	}
}

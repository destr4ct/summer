package rss

import (
	"context"
	"destr4ct/summer/internal/storage"
	"destr4ct/summer/pkg/utils"
	"github.com/mmcdole/gofeed"
)

type Handler struct {
	parser *gofeed.Parser
}

func (h *Handler) ParseSource(ctx context.Context, link string) ([]*storage.SArticle, error) {
	// Получаем служебный формат rss
	feed, err := utils.DoWithAttempts(func() (*gofeed.Feed, error) {
		return h.parser.ParseURLWithContext(link, ctx)
	}, 3)
	if err != nil {
		return nil, err
	}

	articles := make([]*storage.SArticle, 0)
	for _, item := range feed.Items {
		article := storage.NewArticle(item.Link, getBody(item), item.Title)
		articles = append(articles, article)
	}

	return articles, nil
}

func getBody(item *gofeed.Item) string {
	body := item.Content

	for _, plugin := range plugins {
		if body != "" {
			break
		}

		body = plugin(item)
	}

	return body
}

func GetHandler() *Handler {
	return &Handler{
		parser: gofeed.NewParser(),
	}
}

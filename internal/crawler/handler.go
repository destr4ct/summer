package crawler

import (
	"context"
	"destr4ct/summer/internal/crawler/rss"
	"destr4ct/summer/internal/storage"
	"strings"
)

var handlers = map[string]Handler{
	"rss": rss.GetHandler(),
}

type Handler interface {
	ParseSource(ctx context.Context, link string) ([]*storage.SArticle, error)
}

type Source struct {
	Kind string
	Link string
}

func getSource(source string) *Source {
	var result = new(Source)
	assets := strings.Split(source, "|")

	if len(assets) >= 2 {
		result.Kind = assets[0]
		result.Link = assets[1]
	}
	return result
}

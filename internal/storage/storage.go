package storage

import (
	"context"
)

type SummerStorage interface {
	RegisterUser(ctx context.Context, username, tgid string) (*User, error)
	GetUser(ctx context.Context, tgid string) (*User, error)
	AddSource(ctx context.Context, tgid, source string) error
	RemoveSource(ctx context.Context, tgid, source string) error
	AddKeyword(ctx context.Context, tgid, keyword string) error
	RemoveKeyword(ctx context.Context, tgid, keyword string) error

	GetAllSources(ctx context.Context) ([]string, error)
	AddArticle(ctx context.Context, link, content string) (*SArticle, error)

	GetArticleByID(ctx context.Context, id int) (*SArticle, error)
	AddSummary(ctx context.Context, articleID int, summary string) error

	GetArticlesBySource(ctx context.Context, sources string) ([]*SArticle, error)

	Close()
}

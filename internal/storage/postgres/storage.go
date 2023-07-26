package postgres

import (
	"context"
	"destr4ct/summer/internal/config"
	"destr4ct/summer/internal/storage"
	"destr4ct/summer/pkg/utils"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type PgSummerStorage struct {
	pool *pgxpool.Pool
}

func (p *PgSummerStorage) RegisterUser(ctx context.Context, username, tgid string) (*storage.User, error) {
	const insertUserQ = "INSERT INTO summer_user (username, tgid, date_created) VALUES ($1, $2, $3) RETURNING user_id"

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	user := storage.NewUser(username, tgid)
	if err := con.QueryRow(ctx, insertUserQ, user.Username, user.TGID, user.DateCreated).Scan(&user.ID); err != nil {
		if errConverted, ok := err.(*pgconn.PgError); ok && strings.Contains(errConverted.Message, "duplicate") {
			return p.GetUser(ctx, tgid)

		}
		return nil, err
	}
	return user, nil
}

func (p *PgSummerStorage) GetUser(ctx context.Context, tgid string) (*storage.User, error) {
	const getUserQ = `
	SELECT user_id, username, date_created
		FROM summer_user
	WHERE tgid=$1
	`

	const getUserKwQ = "SELECT word FROM keyword WHERE owner_id=$1"
	const getUserSourcesQ = "SELECT link FROM source WHERE owner_id=$1"

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	// Получаем основную информацию
	user := &storage.User{
		TGID:              tgid,
		PreferredKeywords: make([]string, 0),
		Sources:           make([]string, 0),
	}
	if err := con.QueryRow(ctx, getUserQ, tgid).Scan(&user.ID, &user.Username, &user.DateCreated); err != nil {
		return nil, err
	}

	// Запрашиваем источники и ключевые слова из соседних таблиц
	user.PreferredKeywords = p.stringRows(ctx, getUserKwQ, con, user.ID)
	user.Sources = p.stringRows(ctx, getUserSourcesQ, con, user.ID)

	return user, nil
}

func (p *PgSummerStorage) stringRows(ctx context.Context, q string, c *pgxpool.Conn, filter ...interface{}) []string {
	result := make([]string, 0, 8)

	rows, err := c.Query(ctx, q, filter...)
	if err != nil {
		return result
	}

	for rows.Next() {
		var kw string

		if err := rows.Scan(&kw); err != nil {
			break
		}
		result = append(result, kw)
	}
	return result
}

func (p *PgSummerStorage) AddSource(ctx context.Context, tgid, source string) error {
	const addSourceQ = `
		INSERT INTO source(owner_id, link)
		VALUES (
		        (SELECT user_id from summer_user where tgid=$1),
		        $2
		) 
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()

	_, err = con.Exec(ctx, addSourceQ, tgid, source)
	return err
}

func (p *PgSummerStorage) RemoveSource(ctx context.Context, tgid, source string) error {
	const removeSourceQ = `
		DELETE FROM source
		WHERE owner_id=(SELECT user_id FROM summer_user WHERE tgid=$1) AND link=$2
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()

	_, err = con.Exec(ctx, removeSourceQ, tgid, source)
	return err
}

func (p *PgSummerStorage) AddKeyword(ctx context.Context, tgid, keyword string) error {
	const addSourceQ = `
		INSERT INTO keyword(owner_id, word)
		VALUES (
		        (SELECT user_id from summer_user where tgid=$1),
		        $2
		) 
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()

	_, err = con.Exec(ctx, addSourceQ, tgid, keyword)
	return err
}

func (p *PgSummerStorage) RemoveKeyword(ctx context.Context, tgid, keyword string) error {
	const removeSourceQ = `
		DELETE FROM keyword
		WHERE owner_id=(SELECT user_id FROM summer_user WHERE tgid=$1) AND word=$2
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()

	_, err = con.Exec(ctx, removeSourceQ, tgid, keyword)
	return err
}

func (p *PgSummerStorage) GetAllSources(ctx context.Context) ([]string, error) {
	const getSourcesQ = "SELECT DISTINCT link AS src FROM source"

	result := make([]string, 0)

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return result, err
	}
	defer con.Release()

	return p.stringRows(ctx, getSourcesQ, con), nil
}

func (p *PgSummerStorage) AddArticle(ctx context.Context, link, content string, title string) (*storage.SArticle, error) {
	const insertArticleQ = `
		INSERT INTO article(source, content, summary, date_created, has_summary, title)
		VALUES ($1, $2, '', $3, false, $4)
		RETURNING article_id
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	article := storage.NewArticle(link, content, title)
	if err := con.QueryRow(ctx, insertArticleQ, link, content, article.DateCreated, title).Scan(&article.ID); err != nil {
		if errConverted, ok := err.(*pgconn.PgError); ok && strings.Contains(errConverted.Message, "duplicate") {
			return p.GetArticleByTitle(ctx, title)
		}
		return nil, err
	}
	return article, nil
}

func (p *PgSummerStorage) GetArticleByTitle(ctx context.Context, title string) (*storage.SArticle, error) {
	const getPendingQ = `
		SELECT source, content, summary, date_created, has_summary, article_id
			from article
		where title=$1
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	article := storage.SArticle{Title: title}
	err = con.QueryRow(ctx, getPendingQ, article.Title).Scan(
		&article.Source, &article.Content,
		&article.Summary, &article.DateCreated,
		&article.HasSummary, &article.ID,
	)

	return &article, err
}

func (p *PgSummerStorage) GetArticleByID(ctx context.Context, articleID int) (*storage.SArticle, error) {
	const getPendingQ = `
		SELECT source, content, summary, date_created, has_summary, title
			from article
		where article_id=$1
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	article := storage.SArticle{ID: articleID}
	err = con.QueryRow(ctx, getPendingQ, article.ID).Scan(
		&article.Source, &article.Content,
		&article.Summary, &article.DateCreated,
		&article.HasSummary, &article.Title,
	)

	return &article, err
}

func (p *PgSummerStorage) AddSummary(ctx context.Context, articleID int, summary string) error {
	const addSummaryQ = `
		UPDATE article
			SET summary=$1,
			    has_summary=true
		WHERE article_id=$2
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer con.Release()

	_, err = con.Exec(ctx, addSummaryQ, summary, articleID)
	return err
}

func (p *PgSummerStorage) GetArticlesBySource(ctx context.Context, source string) ([]*storage.SArticle, error) {
	const getArticlesBySourceQ = `
		SELECT article_id, source, content, summary, date_created, has_summary, title FROM article
		WHERE source=$1
	`

	con, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Release()

	rows, err := con.Query(ctx, getArticlesBySourceQ, source)
	if err != nil {
		return nil, err
	}

	articles := make([]*storage.SArticle, 8)
	for rows.Next() {
		var na storage.SArticle
		_ = rows.Scan(
			&na.ID, &na.Source, &na.Content,
			&na.Summary, &na.DateCreated,
			&na.HasSummary, &na.Title,
		)
	}
	return articles, nil
}

func (p *PgSummerStorage) Close() {
	p.pool.Close()
}

func GetStorage(cfg *config.DatabaseConfig) (*PgSummerStorage, error) {
	// postgres://username:password@localhost:5432/database_name
	cs := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)

	pool, err := utils.DoWithAttempts(func() (*pgxpool.Pool, error) {
		return pgxpool.New(context.Background(), cs)
	}, 5)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PgSummerStorage{
		pool: pool,
	}, nil
}

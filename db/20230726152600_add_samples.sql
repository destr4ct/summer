-- +goose Up
-- +goose StatementBegin
INSERT INTO summer_user(username, tgid, date_created) values ('cha2ned', 'tg_12345', now());
INSERT INTO source(owner_id, link) values (1, 'rss|https://rssexport.rbc.ru/rbcnews/news/20/full.rss');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

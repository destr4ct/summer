-- +goose Up
-- +goose StatementBegin
CREATE TABLE summer_user (
    user_id SERIAL PRIMARY KEY,
    username varchar(128),
    tgid varchar(32),
    date_created TIMESTAMP
);
CREATE UNIQUE INDEX user_idx ON summer_user (tgid);

CREATE TABLE keyword (
    keyword_id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    word varchar(32),

    FOREIGN KEY (owner_id) REFERENCES summer_user (user_id) ON DELETE CASCADE
);
CREATE TABLE source (
    source_id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    link text,

    FOREIGN KEY (owner_id) REFERENCES summer_user (user_id) ON DELETE CASCADE
);
CREATE TABLE article (
    article_id SERIAL PRIMARY KEY,
    source text,
    content text,
    summary text,
    date_created TIMESTAMP,

    has_summary boolean
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE keyword;
DROP TABLE source;
DROP TABLE summer_user;
DROP TABLE article;
-- +goose StatementEnd

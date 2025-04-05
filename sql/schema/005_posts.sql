-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NULL,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    FOREIGN KEY(feed_id) REFERENCES feeds(id)
);



-- +goose Down
DROP TABLE posts;
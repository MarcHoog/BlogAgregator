-- +goose Up
CREATE TABLE posts (
   id UUID PRIMARY KEY,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
   title TEXT NOT NULL,
   url TEXT NOT NULL,
   description TEXT,
   published_at TIMESTAMP NOT NULL,
   feed_id UUID NOT NULL,
   CONSTRAINT unique_url UNIQUE (url),
   CONSTRAINT fk_feed
       FOREIGN KEY(feed_id)
          REFERENCES feeds(id)
          ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
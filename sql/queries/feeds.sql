-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
       )
    RETURNING *;


-- name: GetFeedByName :one
SELECT * FROM feeds
        WHERE name = $1 LIMIT 1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
        WHERE url = $1 LIMIT 1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
    SET last_fetched_at = $1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;




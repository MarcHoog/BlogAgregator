-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
       )
    RETURNING *;

-- name: GetPostByUrl :one
SELECT * FROM posts
    WHERE url = $1 LIMIT 1;

-- name: GetPostsByFeed :many
SELECT * FROM posts
    WHERE feed_id = $1;

-- name: GetPostsForUser :many
SELECT
    posts.*
FROM posts
    INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
    INNER JOIN feeds ON posts.feed_id = feeds.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.created_at DESC LIMIT $2;


-- name: CheckPostExists :one
SELECT EXISTS (
    SELECT 1 FROM posts WHERE url = $1
);
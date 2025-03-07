// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const checkPostExists = `-- name: CheckPostExists :one
SELECT EXISTS (
    SELECT 1 FROM posts WHERE url = $1
)
`

func (q *Queries) CheckPostExists(ctx context.Context, url string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkPostExists, url)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createPost = `-- name: CreatePost :one
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
    RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostByUrl = `-- name: GetPostByUrl :one
SELECT id, created_at, updated_at, title, url, description, published_at, feed_id FROM posts
    WHERE url = $1 LIMIT 1
`

func (q *Queries) GetPostByUrl(ctx context.Context, url string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByUrl, url)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsByFeed = `-- name: GetPostsByFeed :many
SELECT id, created_at, updated_at, title, url, description, published_at, feed_id FROM posts
    WHERE feed_id = $1
`

func (q *Queries) GetPostsByFeed(ctx context.Context, feedID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByFeed, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT
    posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id
FROM posts
    INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
    INNER JOIN feeds ON posts.feed_id = feeds.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.created_at DESC LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

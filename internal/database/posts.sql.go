// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, description, url, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, feed_id, title, description, url, published_at, created_at, updated_at
`

type CreatePostParams struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	Title       string
	Description sql.NullString
	Url         string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.FeedID,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostByUrl = `-- name: GetPostByUrl :one
SELECT id, feed_id, title, description, url, published_at, created_at, updated_at
FROM posts
where url = $1
`

func (q *Queries) GetPostByUrl(ctx context.Context, url string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByUrl, url)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostsByFeedID = `-- name: GetPostsByFeedID :many
SELECT id, feed_id, title, description, url, published_at, created_at, updated_at
FROM posts
where feed_id = $1
`

func (q *Queries) GetPostsByFeedID(ctx context.Context, feedID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByFeedID, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.PublishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
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
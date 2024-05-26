-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, last_fetched_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT id, name, url, user_id, created_at, updated_at
from feeds;

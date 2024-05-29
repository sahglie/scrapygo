-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, last_fetched_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT id, name, url, user_id, created_at, updated_at
FROM feeds;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetNextFeedsToFetch :many
SELECT *
FROM feeds
ORDER BY last_fetched_at;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = now(),
    updated_at = now()
where id = $1;

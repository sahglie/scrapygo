-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, description, url, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostByUrl :one
SELECT *
FROM posts
where url = $1;

-- name: GetPostsByFeedID :many
SELECT *
FROM posts
where feed_id = $1;

-- name: DeletePostsByFeedID :exec
DELETE
FROM posts
where feed_id = $1;

-- name: GetPostsByUserID :many
SELECT p.*
FROM posts p LEFT JOIN feed_follows ff on ff.feed_id = p.feed_id
where user_id = $1;

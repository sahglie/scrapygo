-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListUsers :many
SELECT id, name, created_at, updated_at from  users;

-- name: FindUserByApiKey :one
SELECT id, name, created_at, updated_at from  users
WHERE api_key = $1;

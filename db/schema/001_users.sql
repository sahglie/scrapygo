-- +goose Up
create table users
(
    id         uuid primary key,
    name       varchar     not null,
    api_key    varchar     not null default encode(sha256(random()::text::bytea), 'hex'),
    created_at timestamptz not null,
    updated_at timestamptz not null
);

-- +goose Down
DROP TABLE users;

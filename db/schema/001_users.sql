-- +goose Up
create table users
(
    id         integer primary key,
    name       varchar not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

-- +goose Down
DROP TABLE users;

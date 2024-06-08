-- +goose Up
create table feeds
(
    id              uuid primary key,
    name            varchar        not null,
    url             varchar unique not null,
    user_id         uuid           not null,
    last_fetched_at timestamptz,
    created_at      timestamptz    not null,
    updated_at      timestamptz    not null,

    foreign key (user_id) references users (id) on delete cascade
);

-- +goose Down
DROP TABLE users;

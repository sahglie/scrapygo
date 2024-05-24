-- +goose Up
create table feed_follows
(
    id         uuid primary key,
    feed_id    uuid        not null,
    user_id    uuid        not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,

    foreign key (user_id) references users(id) on delete cascade,
    foreign key (feed_id) references feeds(id) on delete cascade
);

-- +goose Down
DROP TABLE feed_follows;

-- +goose Up
create table posts
(
    id           uuid primary key,
    feed_id      uuid           not null,
    title        varchar(500)   not null,
    description  text,
    url          varchar unique not null,
    published_at timestamptz    not null,
    created_at   timestamptz    not null,
    updated_at   timestamptz    not null,

    foreign key (feed_id) references feeds (id) on delete cascade
);

-- +goose Down
DROP TABLE posts;

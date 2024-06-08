-- +goose Up
alter table users add constraint unique_name_idx unique(name);

-- +goose Down
alter table users drop constraint unique_name_idx;

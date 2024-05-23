-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR DEFAULT encode(sha256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;



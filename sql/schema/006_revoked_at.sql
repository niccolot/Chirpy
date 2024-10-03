-- +goose Up
ALTER TABLE refresh_tokens
DROP COLUMN revoked_at;

ALTER TABLE refresh_tokens
ADD COLUMN revoked_at TEXT DEFAULT null;

-- +goose Down
ALTER TABLE refresh_tokens
DROP COLUMN revoked_at;

ALTER TABLE refresh_tokens
ADD COLUMN revoked_at TEXT DEFAULT null;


-- +goose Up
ALTER TABLE users
DROP COLUMN is_chirpy_red;

ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL NOT NULL DEFAULT false;



-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;

ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL DEFAULT false;
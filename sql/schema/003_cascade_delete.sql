-- +goose Up
ALTER TABLE chirps
DROP COLUMN user_id;

ALTER TABLE chirps
ADD COLUMN user_id uuid NOT NULL;

ALTER TABLE chirps
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE chirps
DROP CONSTRAINT fk_user;

-- +goose Up
CREATE TABLE refresh_tokens(
    token text primary key not null,
    created_at text not null,
    updated_at text not null,
    user_id uuid unique not null,
    expires_at text not null,
    revoked_at text not null
);

ALTER TABLE refresh_tokens
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;


-- +goose Down
DROP TABLE refresh_tokens;

ALTER TABLE chirps
DROP CONSTRAINT fk_user;

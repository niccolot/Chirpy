-- +goose Up
CREATE TABLE users(
    id uuid primary key not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text not null unique
);

-- +goose Down
DROP TABLE users;

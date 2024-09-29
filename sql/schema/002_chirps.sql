-- +goose Up
CREATE TABLE chirps(
    id uuid primary key not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id uuid not null 
);

-- +goose Down
DROP TABLE chirps;

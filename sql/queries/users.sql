-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: Reset :exec
DELETE FROM users;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

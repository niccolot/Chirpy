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
DELETE FROM users CASCADE;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users 
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpgradeChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;

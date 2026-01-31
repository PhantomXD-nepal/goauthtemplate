-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password_hash
) VALUES (?, ?, ?);

-- name: GetUserByEmail :one
SELECT id, email, password_hash, is_active, created_at
FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, is_active, created_at
FROM users
WHERE id = ?
LIMIT 1;

-- name: DisableUser :exec
UPDATE users
SET is_active = FALSE
WHERE id = ?;

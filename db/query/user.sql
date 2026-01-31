-- name: CreateUser :exec
INSERT INTO users (id, email, password)
VALUES (UUID_TO_BIN(?), ?, ?);

-- name: GetUserByEmail :one
SELECT BIN_TO_UUID(id) as id, email, password, created_at
FROM users
WHERE email = ?;

-- name: GetUserByID :one
SELECT BIN_TO_UUID(id) as id, email, created_at
FROM users
WHERE id = UUID_TO_BIN(?)
LIMIT 1;

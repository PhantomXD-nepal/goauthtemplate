-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    id,
    user_id,
    token_hash,
    expires_at
) VALUES (?, ?, ?, ?);

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, expires_at
FROM refresh_tokens
WHERE id = ?
LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = ?;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = ?;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

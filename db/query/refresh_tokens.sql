-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    id,
    user_id,
    token_hash,
    expires_at
) VALUES (
    UUID_TO_BIN(?),
    UUID_TO_BIN(?),
    ?,
    ?
);

-- name: GetRefreshToken :one
SELECT
    BIN_TO_UUID(id)      AS id,
    BIN_TO_UUID(user_id) AS user_id,
    token_hash,
    expires_at
FROM refresh_tokens
WHERE id = UUID_TO_BIN(?)
LIMIT 1;


-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = UUID_TO_BIN(?);

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

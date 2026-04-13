-- name: CreateRefresh :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 DAY',
    NULL
)
RETURNING *;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

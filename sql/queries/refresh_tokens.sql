-- name: CreateRefreshToken :one

INSERT INTO refresh_tokens (
    token, created_at, updated_at, user_id, expires_at
) VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1
    AND revoked_at IS NULL
    AND expires_at > NOW();

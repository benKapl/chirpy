-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token, 
    created_at, 
    updated_at, 
    user_id, 
    expires_at, 
    revoked_at
)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id AS user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at FROM users
INNER JOIN refresh_tokens
ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
LIMIT 1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;

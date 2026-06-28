-- name: CreateSession :exec
INSERT INTO session(user_id)
VALUES ($1)
RETURNING id;

-- name: GetSessionIDByTokenHash :one
SELECT session_id from access_token
WHERE token_hash = $1;

-- name: CreateAccessToken :one
INSERT INTO access_token(session_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id;

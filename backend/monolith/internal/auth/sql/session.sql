-- name: CreateSession :exec
INSERT INTO session (token_hash, user_id, roles, expires_at)
VALUES ($1, $2, $3, $4);

-- name: DeleteSession :exec
DELETE FROM session WHERE token_hash = $1;

-- name: GetSessionByToken :one
SELECT session_id, user_id, roles, created_at, expires_at
FROM session
WHERE token_hash = $1 AND expires_at > NOW();

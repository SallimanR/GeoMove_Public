-- name: ValidateToken :many
SELECT
	s.user_id,
	s.id AS session_id
FROM access_token AS at
JOIN session AS s ON s.id = sr.session_id
WHERE token_hash = $1
	AND at.expires_at > NOW()
	AND at.revoked_at IS NULL;

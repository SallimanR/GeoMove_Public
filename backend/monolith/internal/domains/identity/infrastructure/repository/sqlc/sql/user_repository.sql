-- name: CreateUser :one
INSERT INTO "user"
DEFAULT VALUES
RETURNING id;

-- name: GetUserRolesByID :many
SELECT role_id
FROM user_roles
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id=$1;

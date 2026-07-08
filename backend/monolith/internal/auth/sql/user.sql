-- name: CreateUser :one
INSERT INTO "user" (phone, email, created_at, updated_at)
VALUES ($1, $2, NOW(), NULL)
RETURNING id;

-- name: UpdateUserPhone :exec
UPDATE "user" SET phone = @phone::TEXT, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE "user" SET email = @email::TEXT, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserProfileImage :exec
UPDATE "user" SET profile_image = @profile_image::TEXT, updated_at = NOW()
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE "user" SET deleted_at = NOW()
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, phone, email, profile_image, created_at, updated_at
FROM "user"
WHERE email = @email::TEXT AND deleted_at IS NULL;

-- name: GetUserByPhone :one
SELECT id, phone, email, profile_image, created_at, updated_at
FROM "user"
WHERE phone = @phone::TEXT AND deleted_at IS NULL;


-- name: CreateOAuthLink :exec
INSERT INTO user_oauth_links (user_id, provider, provider_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, provider) DO UPDATE SET provider_id = EXCLUDED.provider_id;

-- name: GetUserByOAuth :one
SELECT u.id, u.phone, u.email, u.profile_image, u.created_at, u.updated_at
FROM "user" u
JOIN user_oauth_links l ON u.id = l.user_id
WHERE l.provider = $1 AND l.provider_id = $2 AND u.deleted_at IS NULL;

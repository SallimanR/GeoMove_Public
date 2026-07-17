-- name: UpsertSubscription :exec
INSERT INTO push_subscriptions (user_id, endpoint, device_public_key, auth_secret, device_type)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, endpoint)
DO UPDATE SET
    device_public_key = EXCLUDED.device_public_key,
    auth_secret = EXCLUDED.auth_secret,
    device_type = EXCLUDED.device_type,
    created_at = CURRENT_TIMESTAMP;

-- name: GetSubscriptionsByUserID :many
SELECT id, user_id, endpoint, device_public_key, auth_secret, device_type, created_at
FROM push_subscriptions
WHERE user_id = $1;

-- name: DeleteSubscription :exec
DELETE FROM push_subscriptions WHERE endpoint = $1;

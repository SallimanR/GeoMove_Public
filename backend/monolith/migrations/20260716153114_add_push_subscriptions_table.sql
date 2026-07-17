-- migrate:up
CREATE TABLE push_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL,
    device_public_key TEXT NOT NULL,
    auth_secret TEXT NOT NULL,
    device_type TEXT NOT NULL DEFAULT 'web',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, endpoint)
);

CREATE INDEX idx_push_subscriptions_user_id ON push_subscriptions(user_id);

-- migrate:down
DROP TABLE push_subscriptions;

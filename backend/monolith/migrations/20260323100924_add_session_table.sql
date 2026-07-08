-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE session (
    token_hash   TEXT PRIMARY KEY,
    session_id   UUID DEFAULT uuid_generate_v4(),
    user_id      BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMP NOT NULL,
    roles        TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_session_user_id ON session(user_id);
CREATE INDEX idx_session_expires_at ON session(expires_at);

-- migrate:down
DROP TABLE session;

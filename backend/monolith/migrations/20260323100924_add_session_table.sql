-- migrate:up
CREATE TABLE session (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP WITHOUT TIME ZONE,
	ip_adress INET,
	user_agent TEXT
	-- oauth2_id TEXT
);

CREATE TABLE access_token (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	session_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL UNIQUE,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP
);

CREATE INDEX idx_token_hash ON access_token(token_hash);

-- migrate:down
DROP TABLE access_token;
DROP TABLE session;

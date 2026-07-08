-- migrate:up
CREATE TABLE "user" (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE,
	deleted_at TIMESTAMP WITHOUT TIME ZONE, -- soft delete
	phone TEXT UNIQUE,
	email TEXT UNIQUE,
	profile_image TEXT UNIQUE
);

CREATE TABLE user_oauth_links (
    user_id      BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    provider     TEXT NOT NULL,
    provider_id  TEXT NOT NULL,
    PRIMARY KEY (user_id, provider)
);

CREATE INDEX idx_user_oauth_links_provider ON user_oauth_links(provider, provider_id);

-- migrate:down
DROP TABLE user_oauth_links;
DROP TABLE "user";

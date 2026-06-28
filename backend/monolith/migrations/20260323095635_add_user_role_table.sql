-- migrate:up
CREATE TABLE role (
	id INT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE user_roles (
	user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
	role_id INT NOT NULL REFERENCES  role(id) ON DELETE CASCADE,
	assigned_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);

-- migrate:down
DROP TABLE user_role;
DROP TABLE role;

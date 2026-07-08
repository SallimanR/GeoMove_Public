-- migrate:up
ALTER TABLE driver ADD COLUMN user_id BIGINT NOT NULL REFERENCES "user"(id);
CREATE UNIQUE INDEX idx_driver_user_id ON driver(user_id);

-- migrate:down
DROP INDEX IF EXISTS idx_driver_user_id;
ALTER TABLE driver DROP COLUMN user_id;

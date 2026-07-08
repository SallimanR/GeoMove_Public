-- migrate:up
CREATE TABLE driver (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,

	name TEXT NOT NULL,
	profile_image TEXT,

	work_starts TIME,
	work_ends TIME,
	is_available BOOLEAN NOT NULL DEFAULT TRUE,
	last_seen TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	location public.geography(Point, 4326) NOT NULL,
	city_id INTEGER,
	state_id INTEGER,

	rating REAL
);

CREATE INDEX idx_driver_location_geom ON driver USING GIST(location);
CREATE UNIQUE INDEX idx_driver_user_id ON driver(user_id);

-- migrate:down
DROP TABLE driver;

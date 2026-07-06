-- migrate:up
CREATE TABLE driver (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,

	name TEXT NOT NULL,
	profile_image TEXT,

	work_starts TIME WITHOUT TIME ZONE,
	work_ends TIME WITHOUT TIME ZONE,
	is_available BOOLEAN DEFAULT TRUE,
	last_seen TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	location public.geography(Point, 4326) NOT NULL,
	city_id INTEGER,
	state_id INTEGER,

	rating REAL
);

CREATE INDEX idx_driver_location_geom ON driver USING GIST(location);

-- migrate:down
DROP TABLE driver;

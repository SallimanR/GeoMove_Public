-- migrate:up
CREATE TABLE moving_driver (
	driver_id BIGINT PRIMARY KEY REFERENCES driver(user_id) ON DELETE CASCADE,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	realtime_location public.geography(Point, 4326),
	travel_time TIME NOT NULL,
	path_meters INTEGER NOT NULL,
	coarse_h3 H3INDEX GENERATED ALWAYS AS (h3_latlng_to_cell(realtime_location, 2)) STORED,

	destination_location public.geography(Point, 4326),
	destination_time TIMESTAMP
);

CREATE INDEX idx_moving_driver_location ON moving_driver USING GIST(realtime_location);
CREATE INDEX idx_moving_driver_coarse_h3 ON moving_driver(coarse_h3);

CREATE INDEX idx_destination_location ON moving_driver USING GIST(destination_location);
CREATE INDEX idx_destination_time ON moving_driver(destination_time);

-- migrate:down
DROP TABLE moving_driver;

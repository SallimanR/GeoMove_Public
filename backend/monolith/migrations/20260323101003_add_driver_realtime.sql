-- migrate:up
CREATE TABLE driver_realtime (
	driver_id BIGINT PRIMARY KEY REFERENCES driver(id) ON DELETE CASCADE,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	realtime_location public.geography(Point, 4326),
	average_speed REAL,
	predicted_bearing REAL,
	coarse_h3 H3INDEX GENERATED ALWAYS AS (h3_latlng_to_cell(realtime_location, 2)) STORED,

	destination_location public.geography(Point, 4326),
	destination_time TIMESTAMP
);

CREATE INDEX idx_driver_realtime_location ON driver_realtime USING GIST(realtime_location);
CREATE INDEX idx_driver_realtime_coarse_h3 ON driver_realtime(coarse_h3);

CREATE INDEX idx_destination_location ON driver_realtime USING GIST(destination_location);
CREATE INDEX idx_destination_time ON driver_realtime(destination_time);

-- migrate:down
DROP TABLE driver_realtime;

-- migrate:up
CREATE TABLE driver (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- user_id BIGINT NOT NULL UNIQUE REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,

	work_starts TIME WITHOUT TIME ZONE,
	work_ends TIME WITHOUT TIME ZONE,
	is_available BOOLEAN DEFAULT TRUE,
	last_seen TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	location public.geography(Point, 4326) NOT NULL,
	-- city_id integer,
	-- state_id integer,

	-- Equiliteral hexagon
	--                  Area:           Radius:
	h3_res5 H3INDEX, -- ~252,903 km^2
	h3_res6 H3INDEX, -- ~36,129  km^2
	h3_res7 H3INDEX,  -- ~5,161   km^2

	rating REAL
);

CREATE FUNCTION update_driver_location_h3_indexes()
RETURNS TRIGGER AS $$
BEGIN
    NEW.h3_res5 = h3_latlng_to_cell(NEW.location::GEOMETRY, 5);
    NEW.h3_res6 = h3_latlng_to_cell(NEW.location::GEOMETRY, 6);
    NEW.h3_res7 = h3_latlng_to_cell(NEW.location::GEOMETRY, 7);

	NEW.updated_at = CURRENT_TIMESTAMP;

	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_driver_location_update
BEFORE INSERT OR UPDATE OF location ON driver
FOR EACH ROW
EXECUTE FUNCTION update_driver_location_h3_indexes();

-- CREATE INDEX idx_driver_user_id ON driver(user_id);

CREATE INDEX idx_driver_location_geom ON driver USING GIST(location);
CREATE INDEX idx_driver_h3_res5 ON driver(h3_res5);
CREATE INDEX idx_driver_h3_res6 ON driver(h3_res6);
CREATE INDEX idx_driver_h3_res7 ON driver(h3_res7);

-- migrate:down
DROP TRIGGER trigger_driver_location_update;
DROP FUNCTION update_driver_location_h3_indexes;
DROP TABLE driver;

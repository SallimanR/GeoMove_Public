-- migrate:up
CREATE TABLE IF NOT EXISTS country_administrative_boundaries (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,                     -- e.g., 'New York', 'California'
    admin_level integer,                    -- 0: country, 1: state, 2: city, etc.
    geometry geometry(MultiPolygon, 4326) NOT NULL
);

CREATE INDEX IF NOT EXISTS boundaries_geometry_idx 
    ON country_administrative_boundaries USING GIST (geometry);

CREATE INDEX IF NOT EXISTS boundaries_admin_level_idx 
    ON country_administrative_boundaries (admin_level);

-- migrate:down
DROP TABLE country_administrative_boundaries;

-- migrate:up
CREATE TABLE state (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    iso_code varchar(10),
    geometry geometry(MultiPolygon, 4326) NOT NULL,
    bbox geometry(Polygon, 4326),
    h3_res8 h3index[],  -- Pre-computed H3 indexes for faster lookups
    created_at timestamp DEFAULT now()
);

CREATE TABLE city (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    state_id int REFERENCES state(id),
    geometry geometry(MultiPolygon, 4326) NOT NULL,
    bbox geometry(Polygon, 4326),
    h3_res9 h3index[],  -- Higher resolution for cities
    population int,
    timezone text,
    created_at timestamp DEFAULT now()
);

CREATE INDEX idx_state_geometry ON state USING GIST (geometry);
CREATE INDEX idx_state_h3 ON state USING GIN (h3_res8);
CREATE INDEX idx_city_geometry ON city USING GIST (geometry);
CREATE INDEX idx_city_h3 ON city USING GIN (h3_res9);
CREATE INDEX idx_citiy_state ON city (state_id);
-- migrate:down
DROP TABLE city;
DROP TABLE state;

-- migrate:up
CREATE TABLE tow_driver_freely_available (
	user_id BIGINT PRIMARY KEY REFERENCES driver(user_id) ON DELETE CASCADE,
	from_date TIMESTAMP NOT NULL,
	to_date TIMESTAMP NOT NULL,
	from_location public.geography(Point, 4326) NOT NULL,
	from_address TEXT NOT NULL DEFAULT '',
	en_route_order BOOLEAN,
	tariff_per_km REAL
);

CREATE TABLE tow_driver_freely_available_to_location_list (
	id BIGSERIAL PRIMARY KEY,
	tow_driver BIGINT NOT NULL REFERENCES tow_driver_freely_available(user_id) ON DELETE CASCADE,
	location public.geography(Point, 4326) NOT NULL,
	address TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_tow_driv_fa_loc_geom ON tow_driver_freely_available USING GIST(from_location);
CREATE UNIQUE INDEX idx_tow_fa_drv_user_id ON tow_driver_freely_available(user_id);

CREATE INDEX idx_tow_drv_fa_loc_list_geom ON tow_driver_freely_available_to_location_list USING GIST(location);
CREATE INDEX idx_tow_drv_fa_loc_list_fk  ON tow_driver_freely_available_to_location_list(tow_driver);

-- migrate:down
DROP TABLE tow_driver_freely_available_to_location_list;
DROP TABLE tow_driver_freely_available;

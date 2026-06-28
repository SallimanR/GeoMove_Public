-- migrate:up
CREATE TABLE order_scheduled (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	driver_id BIGINT NOT NULL REFERENCES driver(id) ON DELETE CASCADE,
	destination_location NOT NULL public.geography(Point, 4326),
	destination_time TIMESTAMPTZ,
	is_intercity BOOL,
);

CREATE INDEX idx_order_destination_location ON order_scheduled USING GIST(location);
-- migrate:down
DELETE TABLE order_scheduled;

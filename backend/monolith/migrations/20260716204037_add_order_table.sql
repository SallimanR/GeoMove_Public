-- migrate:up
CREATE TYPE ORDER_STATUS AS ENUM (
    'forming',
    'pending',
    'accepted',
    'in_progress',
    'completed',
    'cancelled'
);

CREATE TABLE "order" (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    customer_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    driver_id BIGINT REFERENCES driver(user_id) ON DELETE SET NULL,

    from_location public.geography(Point, 4326) NOT NULL,
    to_location public.geography(Point, 4326) NOT NULL,
	from_address TEXT NOT NULL DEFAULT '',
	to_address TEXT NOT NULL DEFAULT '',
	total_distance_meters INTEGER,
	how_many_wheels_blocked SMALLINT NOT NULL,
    price_rubles INTEGER,

    status ORDER_STATUS NOT NULL DEFAULT 'forming',

    accepted_at TIMESTAMP,
    picked_up_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT
);

CREATE OR REPLACE FUNCTION update_order_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_order_updated_at
    BEFORE UPDATE ON "order"
    FOR EACH ROW
    EXECUTE FUNCTION update_order_updated_at();

CREATE INDEX idx_order_customer_status ON "order"(customer_id, status);
CREATE INDEX idx_order_driver_status ON "order"(driver_id, status);
CREATE INDEX idx_order_active_status ON "order"(status)
    WHERE status IN ('pending', 'accepted', 'in_progress');
CREATE INDEX idx_order_from_location ON "order" USING GIST(from_location);
CREATE INDEX idx_order_to_location ON "order" USING GIST(to_location);

-- migrate:down
DROP TABLE "order" CASCADE;
DROP TYPE ORDER_STATUS;
DROP FUNCTION update_order_updated_at();

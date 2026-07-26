-- migrate:up
CREATE TABLE order_decline (
    order_id BIGINT NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    driver_id BIGINT NOT NULL REFERENCES driver(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (order_id, driver_id)
);

-- migrate:down
DROP TABLE order_decline;

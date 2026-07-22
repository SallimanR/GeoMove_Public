-- migrate:up
CREATE TABLE tow_driver (
    driver_id BIGINT PRIMARY KEY REFERENCES driver(user_id) ON DELETE CASCADE,
    max_car_weight_kg INTEGER NOT NULL,
    max_car_length_meters REAL NOT NULL
);

-- migrate:down
DROP TABLE tow_driver;

-- migrate:up
CREATE TYPE CAR_TYPE AS ENUM (
    'Легковой',
    'Внедорожник',
    'Микроавтобус',
    'Грузовик',
    'Мотоцикл',
    'Спецтехника',
    'Электромобиль',
    'Другое'
);

ALTER TABLE "order"
    ADD COLUMN car_weight_kg INTEGER NOT NULL,
    ADD COLUMN car_length_meters REAL NOT NULL,
    ADD COLUMN car_type CAR_TYPE NOT NULL,
    ADD COLUMN car_name TEXT NOT NULL,
    ADD COLUMN car_photo_url TEXT,
    ADD COLUMN customer_message TEXT;

-- migrate:down
ALTER TABLE "order"
    DROP COLUMN car_weight_kg,
    DROP COLUMN car_length_meters,
    DROP COLUMN car_type,
    DROP COLUMN car_name,
    DROP COLUMN car_photo_url,
    DROP COLUMN customer_message;

DROP TYPE CAR_TYPE;

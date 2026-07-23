-- migrate:up
ALTER TABLE tow_driver ADD COLUMN car_photo_main TEXT NOT NULL DEFAULT '';
ALTER TABLE tow_driver ADD COLUMN car_photos TEXT;

-- migrate:down
ALTER TABLE tow_driver DROP COLUMN IF EXISTS car_photos;
ALTER TABLE tow_driver DROP COLUMN IF EXISTS car_photo_main;

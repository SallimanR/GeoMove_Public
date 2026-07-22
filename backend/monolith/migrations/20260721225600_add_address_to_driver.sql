-- migrate:up
ALTER TABLE driver ADD COLUMN address TEXT;

-- migrate:down
ALTER TABLE driver DROP COLUMN IF EXISTS address;

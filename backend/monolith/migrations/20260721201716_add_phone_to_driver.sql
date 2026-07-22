-- migrate:up
ALTER TABLE driver ADD COLUMN phone TEXT;

-- migrate:down
ALTER TABLE driver DROP COLUMN phone;

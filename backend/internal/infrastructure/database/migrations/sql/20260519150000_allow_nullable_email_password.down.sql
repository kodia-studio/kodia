-- +migrate Down

-- Revert email and password back to NOT NULL (if no NULL values exist)
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;

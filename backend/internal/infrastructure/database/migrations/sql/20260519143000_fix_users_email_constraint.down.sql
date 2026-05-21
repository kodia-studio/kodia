-- +migrate Down

-- Drop the filtered unique index
DROP INDEX IF EXISTS idx_users_email_active;

-- Recreate the original unique constraint
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE(email);

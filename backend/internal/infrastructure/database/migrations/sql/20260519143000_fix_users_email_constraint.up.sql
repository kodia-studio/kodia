-- +migrate Up

-- Drop old unique constraint that includes soft-deleted records
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

-- Create filtered unique index - only enforces uniqueness for active (non-deleted) records
-- This allows multiple soft-deleted records to have the same email
CREATE UNIQUE INDEX CONCURRENTLY idx_users_email_active
ON users(email)
WHERE deleted_at IS NULL;

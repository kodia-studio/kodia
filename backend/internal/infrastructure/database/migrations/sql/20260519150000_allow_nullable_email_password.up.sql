-- +migrate Up

-- Allow email and password to be NULL when account is deleted
-- This enables clearing sensitive data while preserving other user info (name, created_at, etc)
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

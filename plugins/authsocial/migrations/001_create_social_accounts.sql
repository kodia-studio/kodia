-- Up
CREATE TABLE IF NOT EXISTS social_accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider    VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email       VARCHAR(255),
    name        VARCHAR(255),
    avatar_url  TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);
CREATE INDEX IF NOT EXISTS idx_social_accounts_user_id ON social_accounts(user_id);

-- Down
DROP TABLE IF EXISTS social_accounts;

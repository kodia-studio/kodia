-- +migrate Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          VARCHAR(36)  NOT NULL PRIMARY KEY,
    user_id     VARCHAR(36)  NOT NULL,
    token       TEXT         NOT NULL,
    is_revoked  BOOLEAN      NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMP    NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token(255));

-- +migrate Down
DROP TABLE IF EXISTS refresh_tokens;

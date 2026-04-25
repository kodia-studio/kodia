-- +migrate Up

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id          VARCHAR(36)   NOT NULL PRIMARY KEY,
    name        VARCHAR(255)  NOT NULL UNIQUE,
    description TEXT          DEFAULT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP     DEFAULT NULL
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id          VARCHAR(36)   NOT NULL PRIMARY KEY,
    name        VARCHAR(255)  NOT NULL UNIQUE,
    description TEXT          DEFAULT NULL,
    "group"     VARCHAR(100)  NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP     DEFAULT NULL
);

-- Create role_permissions join table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id VARCHAR(36) NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Create user_roles join table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_permissions_group ON permissions("group");
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

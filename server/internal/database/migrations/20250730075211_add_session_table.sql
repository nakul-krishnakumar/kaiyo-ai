-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(128) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    device_info JSONB
);

-- Fast lookup by token hash (most frequent operation)
CREATE UNIQUE INDEX idx_session_hash ON sessions(token_hash);

-- Fast lookup by user (for getting user's active sessions)
CREATE INDEX idx_session_user_id ON sessions(user_id);

-- Composite index for active tokens (not revoked and not expired)
CREATE INDEX idx_session_active ON sessions(user_id, expires_at) 
WHERE revoked_at IS NULL;

-- Fast lookup for expired tokens (for cleanup)
CREATE INDEX idx_session_expires ON sessions(expires_at);

-- Audit index (revoked tokens for cleanup)
CREATE INDEX idx_session_revoked ON sessions(revoked_at) 
WHERE revoked_at IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
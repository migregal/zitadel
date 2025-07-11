CREATE TABLE IF NOT EXISTS auth.pushed_auth_requests (
    instance_id TEXT NOT NULL,
    request_uri TEXT NOT NULL,
    request JSONB NULL,
    expires_at TIMESTAMPTZ NULL,

    PRIMARY KEY (request_uri, instance_id)
);

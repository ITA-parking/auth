CREATE TABLE refresh_token_tbl
(
    id         UUID PRIMARY KEY,
    token      TEXT      NOT NULL UNIQUE,
    user_id    UUID      NOT NULL REFERENCES user_tbl (id),
    expires_at TIMESTAMP NOT NULL,
    revoked    BOOL      NOT NULL DEFAULT false,
    created_at TIMESTAMP          DEFAULT NOW()
);

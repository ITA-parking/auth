CREATE TABLE user_tbl
(
    id         UUID PRIMARY KEY,
    username   TEXT NOT NULL,
    email      TEXT NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE user_tbl
(
    id         UUID PRIMARY KEY,
    username   TEXT NOT NULL,
    email      TEXT NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE identity_tbl
(
    id           UUID PRIMARY KEY,
    pub_key      BYTEA NOT NULL,
    enc_priv_key BYTEA NOT NULL,
    active       bool  NOT NULL,
    user_id      UUID  NOT NULL
);

ALTER TABLE identity_tbl ADD CONSTRAINT fk_identity_user FOREIGN KEY (user_id) REFERENCES user_tbl (id);
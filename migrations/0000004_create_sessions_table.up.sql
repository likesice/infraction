CREATE TABLE sessions
(
    token_hash TEXT PRIMARY KEY NOT NULL,
    user_id    INTEGER          NOT NULL,
    expiry     INTEGER          NOT NULL,
    created_at INTEGER          NOT NULL,
    updated_at INTEGER          NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

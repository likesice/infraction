CREATE TABLE users
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name          TEXT                              NOT NULL,
    email         TEXT UNIQUE                       NOT NULL,
    password_hash TEXT                              NOT NULL,
    created_at    INTEGER                           NOT NULL,
    updated_at    INTEGER                           NOT NULL,
    version       INTEGER                           NOT NULL,
    activated     INTEGER                           NOT NULL DEFAULT 0
);
CREATE TABLE permissions
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    permission  TEXT                              NOT NULL,
    created_at  INTEGER                           NOT NULL,
    updated_at  INTEGER                           NOT NULL,
    version     INTEGER                           NOT NULL
);

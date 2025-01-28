CREATE TABLE infractions
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name       TEXT                              NOT NULL,
    user_id    INTEGER                           NOT NULL,
    created_at INTEGER                           NOT NULL,
    updated_at INTEGER                           NOT NULL,
    version    INTEGER                           NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS secrets (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,

    created_at TEXT NOT NULL,
    expired_at TEXT NOT NULL,

    access_key  TEXT NOT NULL UNIQUE,
    signing_key TEXT NOT NULL,

    message TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS employments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    banner_url TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at TEXT NOT NULL DEFAULT (unixepoch('now')),
    start_date INTEGER NOT NULL,
    end_date INTEGER NOT NULL,
    company TEXT NOT NULL
) STRICT;

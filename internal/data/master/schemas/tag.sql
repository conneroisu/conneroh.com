-- dialect: sqlite
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    icon TEXT NOT NULL DEFAULT 'nf-fa-tag',
    classes TEXT NOT NULL DEFAULT 'uk-badge-secondary',
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch('now'))
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (name);

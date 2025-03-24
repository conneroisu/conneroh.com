CREATE TABLE IF NOT EXISTS embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    embedding TEXT,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now'))
) STRICT;

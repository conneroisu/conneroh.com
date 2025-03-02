-- dialect: sqlite
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT 'nf-fa-tag',
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY (embedding_id) REFERENCES embeddings(id)
);

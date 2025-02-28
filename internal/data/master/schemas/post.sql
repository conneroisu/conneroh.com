-- dialect: sqlite
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    banner_url TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at TEXT NOT NULL DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY(embedding_id) REFERENCES embeddings(id)
);

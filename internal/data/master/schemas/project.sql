-- dialect: sqlite
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    raw_content TEXT NOT NULL,
    banner_url TEXT NOT NULL,
    created_at INTEGER DEFAULT (unixepoch('now')),
    updated_at INTEGER DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY(embedding_id) REFERENCES embeddings(id)
);

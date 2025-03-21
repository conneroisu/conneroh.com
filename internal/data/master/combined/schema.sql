-- Code generated by sqlcquash. DO NOT EDIT.
-- versions: 
--	sqlcquash: v0.0.2


CREATE TABLE IF NOT EXISTS embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    embedding TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now'))
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    raw_content TEXT NOT NULL,
    banner_url TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at TEXT NOT NULL DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY(embedding_id) REFERENCES embeddings(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS post_projects (
    post_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    raw_content TEXT NOT NULL,
    banner_url TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY(embedding_id) REFERENCES embeddings(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS project_tags (
    project_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    raw_content TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT 'tag',
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY (embedding_id) REFERENCES embeddings(id)
) STRICT;

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS project_posts (
    project_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- dialect: sqlite
CREATE TABLE IF NOT EXISTS post_projects (
    post_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
) STRICT;

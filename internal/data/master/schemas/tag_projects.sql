CREATE TABLE IF NOT EXISTS tag_projects (
    tag_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

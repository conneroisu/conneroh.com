CREATE TABLE IF NOT EXISTS work_tags (
    work_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (work_id, tag_id),
    FOREIGN KEY (work_id) REFERENCES employments(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) STRICT;

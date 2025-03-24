CREATE TABLE IF NOT EXISTS employment_tags (
    employment_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (employment_id, tag_id),
    FOREIGN KEY (employment_id) REFERENCES employments(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) STRICT;

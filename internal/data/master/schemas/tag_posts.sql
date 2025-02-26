-- dialect: sqlite
CREATE TABLE IF NOT EXISTS tag_posts (
    tag_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

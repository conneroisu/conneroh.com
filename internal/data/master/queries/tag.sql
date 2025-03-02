-- name: TagGetByID :one
SELECT
    *
FROM
    tags
WHERE
    id = ?
LIMIT
    1;

-- name: TagGetBySlug :one
SELECT
    *
FROM
    tags
WHERE
    slug = ?
LIMIT
    1;

-- name: TagCreate :exec
INSERT INTO
    tags (title, description, slug, embedding_id)
VALUES
    (?, ?, ?, ?);

-- name: TagUpdate :exec
UPDATE
    tags
SET
    title = ?,
    description = ?,
    slug = ?,
    embedding_id = ?
WHERE
    id = ?;

-- name: TagsListByProject :many
SELECT
    t.*
FROM
    tags t
    JOIN project_tags pt ON t.id = pt.tag_id
WHERE
    pt.project_id = ?
ORDER BY
    t.title ASC;

-- name: TagsListByPost :many
SELECT
    t.*
FROM
    tags t
    JOIN post_tags pt ON t.id = pt.tag_id
WHERE
    pt.post_id = ?
ORDER BY
    t.title ASC;

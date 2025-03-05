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
    tags (title, content, raw_content, slug, embedding_id)
VALUES
    (?, ?, ?, ?, ?);

-- name: TagUpdate :exec
UPDATE
    tags
SET
    title = ?,
    slug = ?,
    content = ?,
    icon = ?,
    raw_content = ?,
    embedding_id = ?
WHERE
    id = ?;

-- name: TagsListAlphabetical :many
SELECT
    t.*
FROM
    tags t
ORDER BY
    t.title ASC;

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

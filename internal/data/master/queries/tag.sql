-- name: TagGetByID :one
SELECT
    *
FROM
    tags
WHERE
    id = ?
LIMIT
    1;

-- name: TagGetByName :one
SELECT
    *
FROM
    tags
WHERE
    name = ?
LIMIT
    1;

-- name: TagsListAlphabetical :many
SELECT
    *
FROM
    tags
ORDER BY
    name;

-- name: TagCreate :one
INSERT INTO
    tags (name, description, slug)
VALUES
    (?, ?, ?) RETURNING *;

-- name: TagsListByProject :many
SELECT
    t.*
FROM
    tags t
    JOIN project_tags pt ON t.id = pt.tag_id
WHERE
    pt.project_id = ?
ORDER BY
    t.name ASC;

-- name: TagsListByPost :many
SELECT
    t.*
FROM
    tags t
    JOIN post_tags pt ON t.id = pt.tag_id
WHERE
    pt.post_id = ?
ORDER BY
    t.name ASC;

-- name: PostGet :one
SELECT
    *
FROM
    posts
WHERE
    id = ?
LIMIT
    1;

-- name: PostCreate :one
INSERT INTO
    posts (
        title,
        description,
        slug,
        content,
        raw_content,
        banner_url,
        embedding_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: PostDeleteBySlug :exec
DELETE FROM
    posts
WHERE
    slug = ?;

-- name: PostGetBySlug :one
SELECT
    *
FROM
    posts
WHERE
    slug = ?
LIMIT
    1;

-- name: PostsList :many
SELECT
    *
FROM
    posts
ORDER BY
    created_at DESC;

-- name: PostsListByTag :many
SELECT
    p.*
FROM
    posts p
    JOIN post_tags pt ON p.id = pt.post_id
WHERE
    pt.tag_id = ?
ORDER BY
    p.created_at DESC;

-- name: PostUpdate :one
UPDATE
    posts
SET
    title = ?,
    description = ?,
    slug = ?,
    content = ?,
    raw_content = ?,
    banner_url = ?
WHERE
    id = ? RETURNING *;

-- name: PostDeleteByID :exec
DELETE FROM
    posts
WHERE
    id = ?;

-- name: PostProjectListByPost :many
SELECT
    p.*
FROM
    post_projects pp
    JOIN projects p ON pp.project_id = p.id
WHERE
    pp.post_id = ?
ORDER BY
    p.name ASC;

-- name: PostsListByProject :many
SELECT
    p.*
FROM
    posts p
    JOIN post_projects pp ON p.id = pp.post_id
WHERE
    pp.project_id = ?
ORDER BY
    p.created_at DESC;

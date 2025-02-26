-- name: ProjectGetByID :one
SELECT
    *
FROM
    projects
WHERE
    id = ?
LIMIT
    1;

-- name: ProjectGetBySlug :one
SELECT
    *
FROM
    projects
WHERE
    slug = ?
LIMIT
    1;

-- name: ProjectsList :many
SELECT
    *
FROM
    projects
ORDER BY
    created_at DESC;

-- name: ProjectCreate :one
INSERT INTO
    projects (name, slug, description)
VALUES
    (?, ?, ?) RETURNING *;

-- name: ProjectsListByTag :many
SELECT
    p.*
FROM
    projects p
    JOIN project_tags pt ON p.id = pt.project_id
WHERE
    pt.tag_id = ?
ORDER BY
    p.created_at DESC;

-- name: ProjectUpdate :one
UPDATE
    projects
SET
    name = ?,
    slug = ?,
    description = ?
WHERE
    id = ? RETURNING *;

-- name: ProjectsListByPost :many
SELECT
    p.*
FROM
    projects p
    JOIN project_posts pp ON p.id = pp.project_id
WHERE
    pp.post_id = ?
ORDER BY
    p.created_at DESC;

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

-- name: ProjectCreate :exec
INSERT INTO
    projects (
        title,
        slug,
        description,
        content,
        raw_content,
        banner_url,
        embedding_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?);

-- name: ProjectUpdate :exec
UPDATE
    projects
SET
    title = ?,
    slug = ?,
    description = ?,
    content = ?,
    raw_content = ?,
    banner_url = ?,
    embedding_id = ?
WHERE
    id = ?;

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

-- name: ProjectsListByPost :many
SELECT
    p.*
FROM
    projects p
    JOIN post_projects pp ON p.id = pp.project_id
WHERE
    pp.post_id = ?
ORDER BY
    p.created_at DESC;

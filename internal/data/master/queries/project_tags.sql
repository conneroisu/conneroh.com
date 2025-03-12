-- name: ProjectTagsGetByProjectID :many
SELECT
    *
FROM
    project_tags
WHERE
    project_id = ?;

-- name: ProjectTagsGetByIDs :one
SELECT
    *
FROM
    project_tags
WHERE
    project_id = ?
    AND tag_id = ?;

-- name: ProjectTagCreate :exec
INSERT INTO
    project_tags (project_id, tag_id)
VALUES
    (?, ?);

-- name: ProjectTagDelete :exec
DELETE FROM
    project_tags
WHERE
    project_id = ?
    AND tag_id = ?;

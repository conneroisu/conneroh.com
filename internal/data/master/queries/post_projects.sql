-- name: PostProjectsGetByPostID :many
SELECT
    *
FROM
    post_projects
WHERE
    post_id = ?;

-- name: PostProjectsGetByProjectID :many
SELECT
    *
FROM
    post_projects
WHERE
    project_id = ?;

-- name: PostProjectCreate :exec
INSERT INTO
    post_projects (post_id, project_id)
VALUES
    (?, ?);

-- name: PostProjectDelete :exec
DELETE FROM
    post_projects
WHERE
    post_id = ?
    AND project_id = ?;

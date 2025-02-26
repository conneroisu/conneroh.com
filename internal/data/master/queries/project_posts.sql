-- name: ProjectPostsGetByProjectID :many
SELECT
    *
FROM
    project_posts
WHERE
    project_id = ?;

-- name: ProjectPostsGetByPostID :many
SELECT
    *
FROM
    project_posts
WHERE
    post_id = ?;

-- name: ProjectPostCreate :exec
INSERT INTO
    project_posts (post_id, project_id)
VALUES
    (?, ?);

-- name: ProjectPostDelete :exec
DELETE FROM
    project_posts
WHERE
    project_id = ?
    AND post_id = ?;

-- name: TagPostsGetByTagID :many
SELECT
    *
FROM
    tag_posts
WHERE
    tag_id = ?;

-- name: TagPostsGetByPostID :many
SELECT
    *
FROM
    tag_posts
WHERE
    post_id = ?;

-- name: TagPostsCreate :exec
INSERT INTO
    tag_posts (tag_id, post_id)
VALUES
    (?, ?);

-- name: TagPostsDelete :exec
DELETE FROM
    tag_posts
WHERE
    tag_id = ?
    AND post_id = ?;

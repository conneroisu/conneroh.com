-- name: PostTagsGetByPostID :many
SELECT
    *
FROM
    post_tags
WHERE
    post_id = ?;

-- name: PostTagsGetByTagID :many
SELECT
    *
FROM
    post_tags
WHERE
    tag_id = ?;

-- name: PostTagCreate :exec
INSERT INTO
    post_tags (post_id, tag_id)
VALUES
    (?, ?);

-- name: PostTagDelete :exec
DELETE FROM
    post_tags
WHERE
    post_id = ?
    AND tag_id = ?;

-- name: PostTagGet :one
SELECT
    *
FROM
    post_tags
WHERE
    post_id = ?
    AND tag_id = ?;

-- name: EmploymentTagsList :many
SELECT
    *
FROM
    employment_tags;

-- name: EmploymentTagsGet :one
SELECT
    *
FROM
    employment_tags
WHERE
    employment_id = ?
    AND tag_id = ?;

-- name: EmploymentTagsCreate :exec
INSERT INTO
    employment_tags (employment_id, tag_id)
VALUES
    (?, ?);

-- name: EmploymentTagsDelete :exec
DELETE FROM
    employment_tags
WHERE
    employment_id = ?
    AND tag_id = ?;

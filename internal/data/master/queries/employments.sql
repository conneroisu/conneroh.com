-- name: EmploymentsList :many
SELECT
    *
FROM
    employments
ORDER BY
    created_at DESC;

-- name: EmploymentsGet :one
SELECT
    *
FROM
    employments
WHERE
    id = ?;

-- name: EmploymentsCreate :exec
INSERT INTO
    employments (
        title,
        slug,
        description,
        banner_url,
        start_date,
        end_date
    )
VALUES
    (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    );

-- name: EmploymentsUpdate :exec
UPDATE
    employments
SET
    title = ?,
    slug = ?,
    description = ?,
    banner_url = ?,
    start_date = ?,
    end_date = ?
WHERE
    id = ?;

-- name: EmbeddingsCreate :one
INSERT INTO
    embeddings (embedding, x, y, z)
VALUES
    (?, ?, ?, ?) RETURNING id;

-- name: EmbeddingsGetByID :one
SELECT
    *
FROM
    embeddings
WHERE
    id = ?
LIMIT
    1;

-- name: EmbeddingsUpdate :exec
UPDATE
    embeddings
SET
    embedding = ?,
    x = ?,
    y = ?,
    z = ?
WHERE
    id = ? RETURNING id;

-- name: EmbeddingsCreate :one
INSERT INTO
    embeddings (embedding)
VALUES
    (?) RETURNING id;

-- name: EmbeddingsGetByID :one
SELECT
    id,
    embedding
FROM
    embeddings
WHERE
    id = ?
LIMIT
    1;

-- name: EmeddingUpdate :exec
UPDATE
    embeddings
SET
    embedding = ?
WHERE
    id = ?;

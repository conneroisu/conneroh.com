-- name: GetCompany :one
SELECT * FROM companies
WHERE id = $1;

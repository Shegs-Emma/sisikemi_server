-- name: CreateCollection :one
INSERT INTO collections (
    collection_name
) VALUES (
    $1
)
RETURNING *;

-- name: GetCollection :one
SELECT * FROM collections
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListCollection :many
SELECT * FROM collections
ORDER BY id
LIMIT $1
OFFSET $2;
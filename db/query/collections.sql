-- name: CreateCollection :one
INSERT INTO collections (
    collection_name,
    collection_description,
    product_count,
    thumbnail_image,
    header_image
) VALUES (
    $1, $2, $3, $4, $5
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

-- name: UpdateCollection :one
UPDATE collections
SET
  collection_name = COALESCE(sqlc.narg(collection_name), collection_name),
  collection_description = COALESCE(sqlc.narg(collection_description), collection_description),
  product_count = COALESCE(sqlc.narg(product_count), product_count),
  thumbnail_image = COALESCE(sqlc.narg(thumbnail_image), thumbnail_image),
  header_image = COALESCE(sqlc.narg(header_image), header_image)
WHERE
  id = sqlc.arg(id)
RETURNING *;
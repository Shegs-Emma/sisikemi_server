-- name: CreateProductMedia :one
INSERT INTO product_media (
    product_media_ref,
    product_id,
    is_main_image,
    media_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProductMedia :one
SELECT * FROM product_media
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetProductMediaByRef :one
SELECT * FROM product_media
WHERE product_media_ref = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListProductMedia :many
SELECT * FROM product_media
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetProductsMedium :many
SELECT * FROM product_media
WHERE product_id = $1
ORDER BY id
LIMIT $2;

-- name: DeleteProductMedia :exec
DELETE FROM product_media
WHERE product_id = $1;
-- name: CreateCartItem :one
INSERT INTO cart (
    product_id, 
    product_name,
    user_ref_id,
    product_price, 
    product_quantity,
    product_image,
    product_color,
    product_size
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetCartItem :one
SELECT * FROM cart
WHERE id = $1 LIMIT 1;

-- name: GetCartItemByProductId :one
SELECT * FROM cart
WHERE product_id = $1 LIMIT 1;

-- name: GetCartItemByUser :one
SELECT * FROM cart
WHERE user_ref_id = $1 LIMIT 1;

-- name: ListCartItems :many
SELECT * FROM cart
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListUserCartItems :many
SELECT * FROM cart
WHERE user_ref_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateCartItemQty :one
UPDATE cart
SET
    product_quantity = COALESCE(sqlc.narg(product_quantity), product_quantity)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart
WHERE product_id = $1;
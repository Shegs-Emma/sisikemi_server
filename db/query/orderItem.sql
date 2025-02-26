-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListOrderItems :many
SELECT * FROM order_items
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetOrderItemsForOrder :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY id
LIMIT $2;
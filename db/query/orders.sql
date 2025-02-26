-- name: CreateOrder :one
INSERT INTO orders (
    ref_no,
    username,
    amount,
    payment_method,
    order_status,
    shipping_method,
    shipping_address_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetUserOrders :many
SELECT * FROM orders
WHERE username = $1
ORDER BY id
LIMIT $2;
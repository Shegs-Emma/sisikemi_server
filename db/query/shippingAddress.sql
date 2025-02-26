-- name: CreateShippingAddress :one
INSERT INTO shipping_address (
    username,
    country,
    address,
    town,
    postal_code,
    landmark
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetShippingAddress :one
SELECT * FROM shipping_address
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListShippingAddresses :many
SELECT * FROM shipping_address
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetUserShippingAddresses :many
SELECT * FROM shipping_address
WHERE username = $1
ORDER BY id
LIMIT $2;
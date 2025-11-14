-- name: CreateProduct :one
INSERT INTO products (
    product_ref_no,
    product_name,
    product_description,
    product_code,
    price,
    sale_price,
    product_image_main,
    product_image_other_1,
    product_image_other_2,
    product_image_other_3,
    collection,
    quantity,
    color,
    size,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET
  product_ref_no = COALESCE(sqlc.narg(product_ref_no), product_ref_no),
  product_name = COALESCE(sqlc.narg(product_name), product_name),
  product_description = COALESCE(sqlc.narg(product_description), product_description),
  product_code = COALESCE(sqlc.narg(product_code), product_code),
  price = COALESCE(sqlc.narg(price), price),
  sale_price = COALESCE(sqlc.narg(sale_price), sale_price),
  product_image_main = COALESCE(sqlc.narg(product_image_main), product_image_main),
  product_image_other_1 = COALESCE(sqlc.narg(product_image_other_1), product_image_other_1),
  product_image_other_2 = COALESCE(sqlc.narg(product_image_other_2), product_image_other_2),
   product_image_other_3 = COALESCE(sqlc.narg(product_image_other_3), product_image_other_3),
  collection = COALESCE(sqlc.narg(collection), collection),
  quantity = COALESCE(sqlc.narg(quantity), quantity),
  color = COALESCE(sqlc.narg(color), color),
  size = COALESCE(sqlc.narg(size), size),
  status = COALESCE(sqlc.narg(status), status)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: CountProducts :one
SELECT COUNT(*) AS total FROM products;

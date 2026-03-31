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
SELECT 
    p.id,
    p.product_ref_no,
    p.product_name,
    p.product_description,
    p.product_code,
    p.price,
    p.sale_price,
    p.quantity,
    p.color,
    p.size,
    p.status,
    p.created_at,
    p.collection,
    p.product_image_main,
    p.product_image_other_1,
    p.product_image_other_2,
    p.product_image_other_3,

    -- Collection
    c.collection_name AS collection_name,

    -- Image URLs
    m_main.url AS main_image_url,
    m_o1.url AS other_image_1_url,
    m_o2.url AS other_image_2_url,
    m_o3.url AS other_image_3_url

FROM products p
LEFT JOIN collections c ON p.collection = c.id
LEFT JOIN product_media pm_main 
    ON pm_main.product_media_ref = p.product_image_main 
LEFT JOIN media m_main 
    ON m_main.media_ref = pm_main.media_id
LEFT JOIN product_media pm_o1   
    ON pm_o1.product_media_ref = p.product_image_other_1 
LEFT JOIN media m_o1   
    ON m_o1.media_ref = pm_o1.media_id
LEFT JOIN product_media pm_o2   
    ON pm_o2.product_media_ref = p.product_image_other_2 
LEFT JOIN media m_o2   
    ON m_o2.media_ref = pm_o2.media_id
LEFT JOIN product_media pm_o3   
    ON pm_o3.product_media_ref = p.product_image_other_3 
LEFT JOIN media m_o3   
    ON m_o3.media_ref = pm_o3.media_id

WHERE
    (sqlc.narg('search')::text IS NULL OR p.product_name ILIKE '%' || sqlc.narg('search')::text || '%')
    AND (sqlc.narg('collection')::bigint IS NULL OR p.collection = sqlc.narg('collection')::bigint)
    AND (sqlc.narg('min_price')::bigint IS NULL OR p.price >= sqlc.narg('min_price')::bigint)
    AND (sqlc.narg('max_price')::bigint IS NULL OR p.price <= sqlc.narg('max_price')::bigint)
    AND (sqlc.narg('product_name')::text IS NULL OR p.product_name ILIKE '%' || sqlc.narg('product_name')::text || '%')

ORDER BY
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'price' AND sqlc.narg('sort_order')::text = 'asc' 
        THEN p.price 
    END ASC,
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'price' AND sqlc.narg('sort_order')::text = 'desc' 
        THEN p.price 
    END DESC,
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'created_at' AND sqlc.narg('sort_order')::text = 'asc' 
        THEN p.created_at 
    END ASC,
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'created_at' AND sqlc.narg('sort_order')::text = 'desc' 
        THEN p.created_at 
    END DESC,
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'product_name' AND sqlc.narg('sort_order')::text = 'asc' 
        THEN p.product_name 
    END ASC,
    CASE 
        WHEN sqlc.narg('sort_field')::text = 'product_name' AND sqlc.narg('sort_order')::text = 'desc' 
        THEN p.product_name 
    END DESC,
    p.id DESC

LIMIT $1 OFFSET $2;

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
SELECT COUNT(*)
FROM products
WHERE
    (sqlc.narg(search)::text IS NULL OR product_name ILIKE '%' || sqlc.narg(search) || '%')
AND
    (sqlc.narg(collection)::bigint IS NULL OR collection = sqlc.narg(collection)::bigint)
AND
    (sqlc.narg(min_price)::bigint IS NULL OR price >= sqlc.narg(min_price))
AND
     (sqlc.narg(max_price)::bigint IS NULL OR price <= sqlc.narg(max_price))
AND
    (sqlc.narg(product_name)::text IS NULL OR product_name ILIKE '%' || sqlc.narg(product_name) || '%');

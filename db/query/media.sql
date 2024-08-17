-- name: CreateMedia :one
INSERT INTO media (
    media_ref, url, aws_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMedia :one
SELECT * FROM media
WHERE media_ref = $1 LIMIT 1;

-- name: GetMediaForUpdate :one
SELECT * FROM media
WHERE media_ref = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListMedia :many
SELECT * FROM media
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateMedia :one
UPDATE media
SET url = $2, aws_id = $3
WHERE media_ref = $1
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media
WHERE media_ref = $1;
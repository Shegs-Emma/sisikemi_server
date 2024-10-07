-- name: CreateMedia :one
INSERT INTO media (
    media_ref, url, aws_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMedia :one
SELECT * FROM media
WHERE id = $1 LIMIT 1;

-- name: GetMediaByRef :one
SELECT * FROM media
WHERE media_ref = $1 LIMIT 1;

-- name: ListMedia :many
SELECT * FROM media
ORDER BY id
LIMIT $1
OFFSET $2;
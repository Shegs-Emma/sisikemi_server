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
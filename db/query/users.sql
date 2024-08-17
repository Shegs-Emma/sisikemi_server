-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  first_name,
  last_name,
  phone_number,
  profile_photo,
  email,
  is_admin
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUsers :one
UPDATE users
SET username = $2, first_name = $3, last_name = $4, phone_number = $5, profile_photo = $6, email = $7, is_admin = $8, hashed_password = $9
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
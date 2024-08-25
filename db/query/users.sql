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

-- -- name: GetUserForUpdate :one
-- SELECT * FROM users
-- WHERE username = $1 LIMIT 1
-- FOR NO KEY UPDATE;

-- -- name: ListUsers :many
-- SELECT * FROM users
-- ORDER BY username
-- LIMIT $1
-- OFFSET $2;

-- -- name: UpdateUsers :one
-- UPDATE users
-- SET first_name = $2, last_name = $3, phone_number = $4, profile_photo = $5, email = $6, is_admin = $7, hashed_password = $8
-- WHERE username = $1
-- RETURNING *;

-- -- name: DeleteUser :exec
-- DELETE FROM users
-- WHERE username = $1;
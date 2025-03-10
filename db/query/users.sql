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
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  profile_photo = COALESCE(sqlc.narg(profile_photo), profile_photo)
WHERE
  username = sqlc.arg(username)
RETURNING *;

-- name: UpdateUserVerificationCode :one
UPDATE users
SET
  verification_code = COALESCE(sqlc.narg(verification_code), verification_code)
WHERE
  email = sqlc.arg(email)
RETURNING *;
  

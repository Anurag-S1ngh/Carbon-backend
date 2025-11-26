-- name: CreateUser :one
INSERT INTO users (
  email,
  profile_image_url
) VALUES ( $1, $2 ) 
RETURNING *;

-- name: GetUserByEmail :one
SELECT email, profile_image_url
  FROM users
  WHERE email = $1 AND deleted_at IS NULL;


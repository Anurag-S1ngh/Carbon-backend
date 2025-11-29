-- name: CreateUser :one
INSERT INTO users (
  email,
  profile_image_url,
  github_username,
  github_access_token
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
  WHERE email = $1 AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET
  email               = COALESCE($1, email),
  profile_image_url   = COALESCE($2, profile_image_url),
  github_username     = COALESCE($3, github_username),
  github_access_token = COALESCE($4, github_access_token),
  updated_at = NOW()
WHERE email = $1
RETURNING *;
